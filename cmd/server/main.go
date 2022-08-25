package main

import (
	"log"

	"github.com/ShauryaAg/opfs/cmd"
	"github.com/ShauryaAg/opfs/service"
)

func main() {
	server := service.NewNodeServer(cmd.Name, cmd.Addr, cmd.RoutingTable)
	log.Fatal(server.StartListening())
}
