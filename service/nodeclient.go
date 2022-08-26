package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/ShauryaAg/opfs/cmd"
	"github.com/ShauryaAg/opfs/config"
	"github.com/ShauryaAg/opfs/pb/pb"
	"github.com/ShauryaAg/opfs/types"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeClient struct {
	name string
	addr net.TCPAddr

	Files []types.File // all the files that the user has

	wg sync.WaitGroup

	routingTable *types.RoutingTable

	Peers      map[string]types.Node
	PeerClient map[string]pb.UploadFileServiceClient

	PeerChunk map[string]string     // to store which oeer has what chunk, that is, {peer_id -> chunk_id}
	ChunkPeer map[string]types.Node // to store which chunk is stored with which peer, that is, {chunk_id -> peer_id}
}

func NewNodeClient(name string, addr net.TCPAddr, routingTable *types.RoutingTable) *NodeClient {
	bootstrappeers := config.StartPeers
	peers := make(map[string]types.Node)

	for _, node := range bootstrappeers {
		id := uuid.New().String()
		peers[id] = types.Node{Name: id, Ip: node.IP.String(), Port: int32(node.Port)}
	}

	peerChunk := make(map[string]string)
	chunkPeer := make(map[string]types.Node)
	return &NodeClient{name: name, addr: addr, Peers: peers, routingTable: routingTable, PeerChunk: peerChunk, ChunkPeer: chunkPeer}
}

// Upload the chunk to all the known peers
func (nc *NodeClient) UploadFile(ctx context.Context, in *pb.UploadFileRequest, opts ...grpc.CallOption) (*pb.UploadFileResponse, error) {
	ch := make(chan pb.UploadFileResponse, len(nc.Peers))
	print("gergloer")

	for _, peer := range nc.Peers {
		fmt.Println(peer)
		go func(peer types.Node) error {
			loc := fmt.Sprintf("%s:%d", peer.Ip, peer.Port)
			conn, err := grpc.Dial(loc, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("encountered error while connecting to: %s: \n%v", loc, err)
				return err
			}

			client := pb.NewUploadFileServiceClient(conn)
			response, err := client.UploadFile(ctx, in)
			if err != nil {
				response = &pb.UploadFileResponse{Ack: false}
			}

			nc.PeerChunk[peer.Name] = in.GetChunk().Id
			nc.ChunkPeer[in.GetChunk().Id] = peer

			fmt.Println(response.String())
			ch <- *response
			return nil
		}(peer)
	}

	response := <-ch
	// response := pb.UploadFileResponse{Ack: true}
	return &response, nil
}

func (nc *NodeClient) ShareFileWithNode(ctx context.Context, file types.File, node types.Node) (bool, error) {
	loc := fmt.Sprintf("%s:%d", node.Ip, node.Port)
	conn, err := grpc.Dial(loc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server for node: %s: %v", loc, err)
		return false, err
	}

	var chunkdata []*pb.ChunkData
	for _, id := range file.Sequence {
		pbnode := pb.Node{Name: nc.ChunkPeer[id].Name, Ip: nc.ChunkPeer[id].Ip, Port: nc.ChunkPeer[id].Port}
		chunkdata = append(chunkdata, &pb.ChunkData{Node: &pbnode, Chunkid: id})
	}

	req := pb.ShareFileDataRequest{Chunkdata: chunkdata, Sequence: file.Sequence}
	client := pb.NewUploadFileServiceClient(conn)
	response, err := client.ShareFileData(ctx, &req)
	if err != nil {
		log.Printf("error while calling ShareFileData method %v", err)
		return false, err
	}

	return response.Ack, nil
}

func (nc *NodeClient) UploadFileToServer(ctx context.Context, file types.File) {
	for index, id := range file.Sequence {
		nc.wg.Add(1)

		go func(chunk types.Chunk, index int) {
			defer nc.wg.Done()

			log.Printf("uploading chunk : %d / %d", index+1, len(file.Chunks))
			chunkdata := pb.Chunk{Id: chunk.Id, Data: chunk.Data}
			req := pb.UploadFileRequest{Chunk: &chunkdata, Node: &pb.Node{Name: cmd.Name, Ip: cmd.Addr.IP.String(), Port: int32(cmd.Addr.Port)}}
			nc.UploadFile(ctx, &req)
		}(file.Chunks[id], index)
	}

	log.Print("wating...")
	nc.wg.Wait()
	log.Printf("Uploaded file to server")
}
