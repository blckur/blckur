package shard

import (
	"github.com/blckur/blckur/utils"
	"hash/fnv"
	"sort"
)

type Shard struct {
	defaultConsistency int
	maxConsistency int
	servers [][][]string
}

func (s *Shard) SelectX(key string, consistency int) (servers []string) {
	n := len(s.servers)
	if n == 0 {
		servers = []string{}
		return
	}
	consistency = utils.MinInt(consistency, s.maxConsistency)

	hash := fnv.New32a()
	hash.Write([]byte(key))
	servers = s.servers[consistency - 1][hash.Sum32() % uint32(n)]

	return
}

func (s *Shard) SelectAll(key string) (servers []string) {
	n := len(s.servers)
	if n == 0 {
		servers = []string{}
		return
	}

	hash := fnv.New32a()
	hash.Write([]byte(key))
	servers = s.servers[s.maxConsistency][hash.Sum32() % uint32(n)]

	return
}

func (s *Shard) Select(key string) (servers []string) {
	n := len(s.servers)
	if n == 0 {
		servers = []string{}
		return
	}

	hash := fnv.New32a()
	hash.Write([]byte(key))
	servers = s.servers[s.defaultConsistency - 1][hash.Sum32() % uint32(n)]

	return
}

func New(keys []string, defaultConsistency int) (shd *Shard) {
	maxConsistency := len(keys)
	defaultConsistency = utils.MinInt(defaultConsistency, maxConsistency)

	shd = &Shard{
		defaultConsistency: defaultConsistency,
		maxConsistency: maxConsistency,
		servers: make([][][]string, len(keys)),
	}

	sort.Strings(keys)

	for k := 0; k < maxConsistency; k++ {
		servers := make([][]string, len(keys))
		for j := 0; j < k + 1; j++ {
			if j != 0 {
				utils.RotateStrings(keys, 1)
			}
			for i, key := range keys {
				servers[i] = append(servers[i], key)
			}
		}
		shd.servers[k] = servers
	}

	return
}
