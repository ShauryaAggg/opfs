package types

import "github.com/golang/groupcache/lru"

type RoutingTable struct {
	Routes *lru.Cache
}

func NewRoutingTable(maxNeighbours int) *RoutingTable {
	return &RoutingTable{
		Routes: lru.New(maxNeighbours),
	}
}
