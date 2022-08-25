package cmd

import (
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
	Name = uuid.New().String()
	Addr = net.TCPAddr{IP: []byte{0, 0, 0, 0}, Port: 8080}
	RoutingTable = types.NewRoutingTable(config.MaxNeighbours)
}
