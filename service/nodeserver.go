package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/ShauryaAg/opfs/config"
	pb "github.com/ShauryaAg/opfs/pb/pb"
	"github.com/ShauryaAg/opfs/types"
	"github.com/ShauryaAg/opfs/utils"
)

type NodeServer struct {
	name string
	addr net.TCPAddr

	routingTable *types.RoutingTable

	Peers map[string]pb.UploadFileServiceClient // this is what stores the peers (routing table basically)
	Files []types.File

	PeerChunk map[string]string // to store which oeer has what chunk, that is, {peer_id -> chunk_id}
	ChunkPeer map[string]string // to store which chunk is stored with which peer, that is, {chunk_id -> peer_id}

	Store map[string][]byte // list of chunks that are stored with the NodeServer
}

func NewNodeServer(name string, addr net.TCPAddr, routingTable *types.RoutingTable) *NodeServer {
	store := make(map[string][]byte)
	return &NodeServer{name: name, addr: addr, routingTable: routingTable, Store: store}
}

func (ns *NodeServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	pbnode := req.GetNode()
	pbchunk := req.GetChunk()

	ns.routingTable.Routes.Add(pbnode.GetName(), types.Node{pbnode.Name, pbnode.Ip, pbnode.Port})
	// also check if the storing was successful or not and return Ack based on this
	ns.Store[pbchunk.GetId()] = pbchunk.GetData()

	fmt.Println(ns.Store)
	res := &pb.UploadFileResponse{Ack: true}
	return res, nil
}

func (ns *NodeServer) ShareFileData(ctx context.Context, req *pb.ShareFileDataRequest) (*pb.ShareFileDataResponse, error) {
	chunkdata := req.GetChunkdata()
	filemap := make(map[string][]byte)

	for _, data := range chunkdata {
		pbnode := data.GetNode()
		pbchunkid := data.GetChunkid()

		node := types.Node{Name: pbnode.Name, Ip: pbnode.Ip, Port: pbnode.Port}
		chunk, err := ns.FetchChunkFromNode(ctx, node, pbchunkid)
		if err != nil {
			log.Printf("error fetching data from node: %v: for chunkid: %s", node, pbchunkid)
			continue
		}

		filemap[pbchunkid] = chunk
	}

	var chunks []types.Chunk
	for id, data := range filemap {
		chunks = append(chunks, types.Chunk{Id: id, Data: data})
	}

	file := new(types.File)
	file.Chunks = chunks
	binary := file.JoinChunks()
	ehr, err := utils.ConvertBinaryToEhr(binary)
	if err != nil {
		return &pb.ShareFileDataResponse{Ack: false}, status.Errorf(codes.Internal, "couldn't create file")
	}

	fmt.Println(ehr) // complete ehr

	resp := &pb.ShareFileDataResponse{Ack: true}
	return resp, nil
}

func (ns *NodeServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	chunkid := req.GetChunkid()
	data, ok := ns.Store[chunkid]
	if !ok {
		log.Printf("store doesn't contain chunkid: %s", chunkid)
		response := &pb.DownloadFileResponse{}
		return response, nil
	}

	chunk := pb.Chunk{Id: chunkid, Data: data}
	response := &pb.DownloadFileResponse{Chunk: &chunk}
	return response, nil
}

func (ns *NodeServer) FetchChunkFromNode(ctx context.Context, node types.Node, chunkid string) ([]byte, error) {
	loc := fmt.Sprintf("%s:%d", node.Ip, node.Port)
	conn, err := grpc.Dial(loc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to client: %s", loc)
		return nil, err
	}

	client := pb.NewUploadFileServiceClient(conn)

	req := &pb.DownloadFileRequest{Node: &pb.Node{Name: node.Name, Ip: node.Ip, Port: node.Port}, Chunkid: chunkid}
	response, err := client.DownloadFile(ctx, req)
	if err != nil {
		log.Printf("error downloading chunk: %s from node: %s", chunkid, loc)
		return nil, err
	}

	chunk := response.GetChunk()
	return chunk.GetData(), nil
}

func (ns *NodeServer) StartListening() error {
	address := fmt.Sprintf("%s:%d", ns.addr.IP.String(), ns.addr.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Couldn't start listening: %v", err)
		return err
	}

	server := grpc.NewServer()
	pb.RegisterUploadFileServiceServer(server, ns)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Printf("Error while starting server %v", err)
		return err
	}

	log.Printf("starting grpc server on: %d", ns.addr.Port)
	return nil
}

func (ns *NodeServer) pingHello(addr net.TCPAddr) error {
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}

	conn.Write([]byte(ns.name))
	return nil
}

// This method is ran when it first joins the network to inform other NodeServers about itself
func (ns *NodeServer) Bootstrap() {
	for _, peer := range config.StartPeers {
		// make request and update your peers
		ns.pingHello(peer)
	}
}

func (ns *NodeServer) Setup(name string, addr net.TCPAddr) error {
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error setting up NodeServer")
		return err
	}
	defer conn.Close()

	ns.Peers[name] = pb.NewUploadFileServiceClient(conn)

	return nil
}
