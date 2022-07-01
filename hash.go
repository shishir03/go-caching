package main

import (
	"github.com/buraksezer/consistent"
	"hash/fnv"
)

type hasher struct{}

func (hash hasher) Sum64(s []byte) uint64 {
	h := fnv.New64a()
	h.Write(s)
	return h.Sum64()
}

func ConsistentHash(n int) *consistent.Consistent {
	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    n,
		ReplicationFactor: 20,
		Load:              1.25,
	}

	return consistent.New(nil, cfg)
}
