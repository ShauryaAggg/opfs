package main

import (
	"context"

	"github.com/ShauryaAg/opfs/cmd"
	"github.com/ShauryaAg/opfs/config"
	"github.com/ShauryaAg/opfs/service"
	"github.com/ShauryaAg/opfs/types"
	"github.com/ShauryaAg/opfs/utils"
)

func uploadFile(client *service.NodeClient) error {
	ehr := types.Ehr{Id: "123", Patient: "Name", Details: []types.EhrDetail{
		{Id: "1244", Date: "1/12/99", Description: "description", Type: "sometype", Value: "value"},
	}}

	data, err := utils.ConvertEhrToBinary(ehr)
	if err != nil {
		return err
	}

	chunks := utils.DivideIntoChunks(data, config.ChunkSize)
	file := types.NewFile("file-one", chunks)

	client.UploadFileToServer(context.Background(), *file)
	// client.ShareFileWithNode(context.Background(), file, node)

	return nil
}

func main() {
	client := service.NewNodeClient(cmd.Name, cmd.Addr, cmd.RoutingTable)

	uploadFile(client)
}
