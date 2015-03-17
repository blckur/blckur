package shard

import (
	"github.com/blckur/blckur/utils"
	"hash/fnv"
	"sort"
)

type Shard struct {
	servers [][]string
}

func (s *Shard) Select(key string) (servers []string) {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	servers = s.servers[hash.Sum32() % uint32(len(s.servers))]
	return
}

func New(keys []string, consistency int) (shd *Shard) {
	consistency = utils.MinInt(consistency, len(keys))
	shd = &Shard{
		servers: make([][]string, len(keys)),
	}

	sort.Strings(keys)

	for j := 0; j < consistency; j++ {
		if j != 0 {
			utils.RotateStrings(keys, 1)
		}
		for i, key := range keys {
			shd.servers[i] = append(shd.servers[i], key)
		}
	}

	return
}
