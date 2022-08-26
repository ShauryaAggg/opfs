package cmd

import (
	"flag"
	"net"

	"github.com/ShauryaAg/opfs/config"
	"github.com/ShauryaAg/opfs/types"
	"github.com/google/uuid"
)

var (
	Name         string
	Addr         net.TCPAddr
	RoutingTable *types.RoutingTable
)

func init() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	Name = uuid.New().String()
	Addr = net.TCPAddr{IP: []byte{0, 0, 0, 0}, Port: *port}
	RoutingTable = types.NewRoutingTable(config.MaxNeighbours)
}
