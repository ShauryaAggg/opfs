package config

import (
	"net"
)

var MaxNeighbours int = 20
var RedundancyFactor int = 3

var StartPeers []net.TCPAddr = []net.TCPAddr{
	{IP: []byte{0, 0, 0, 0}, Port: 8080},
}

var ChunkSize = 16
