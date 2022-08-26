package main

import (
	"context"
	"fmt"

	"github.com/ShauryaAg/opfs/cmd"
	"github.com/ShauryaAg/opfs/config"
	"github.com/ShauryaAg/opfs/service"
	"github.com/ShauryaAg/opfs/types"
	"github.com/ShauryaAg/opfs/utils"
	"github.com/google/uuid"
)

func uploadFile(client *service.NodeClient) error {
	ehr := types.Ehr{Id: "123", Patient: "Name", Details: []types.EhrDetail{
		{Id: "1244", Date: "1/12/99", Description: "description", Type: "sometype", Value: "value"},
	}}

	data, err := utils.ConvertEhrToBinary(ehr)
	if err != nil {
		return err
	}

	chunks, sequence := utils.DivideIntoChunks(data, config.ChunkSize)
	fmt.Printf("BEFORE: %v", chunks)
	file := types.NewFile("file-one", chunks, sequence)

	client.UploadFileToServer(context.Background(), *file)

	node := types.Node{Name: uuid.New().String(), Ip: "0.0.0.0", Port: 8081}
	client.ShareFileWithNode(context.Background(), *file, node)

	return nil
}

func main() {
	client := service.NewNodeClient(cmd.Name, cmd.Addr, cmd.RoutingTable)

	uploadFile(client)
}
