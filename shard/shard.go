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

func New(servers map[string]string, consistency int) (shd *Shard) {
	consistency = utils.MinInt(consistency, len(servers))
	keys := make([]string, len(servers))
	shd = &Shard{
		servers: make([][]string, len(servers)),
	}

	i := 0
	for key, _ := range servers {
		keys[i] = key
		i += 1
	}

	sort.Strings(keys)

	for j := 0; j < consistency; j++ {
		if j != 0 {
			utils.RotateStrings(keys, 1)
		}
		for i, key := range keys {
			shd.servers[i] = append(shd.servers[i], servers[key])
		}
	}

	return
}
