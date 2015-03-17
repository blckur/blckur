package distribute

import (
	"github.com/blckur/blckur/utils"
	"hash/fnv"
	"sort"
)

type Distribute struct {
	servers [][]string
}

func (d *Distribute) Select(key string) (servers []string) {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	servers = d.servers[hash.Sum32() % uint32(len(d.servers))]
	return
}

func New(servers map[string]string, consistency int) (dist *Distribute) {
	consistency = utils.MinInt(consistency, len(servers))
	keys := make([]string, len(servers))
	dist = &Distribute{
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
			dist.servers[i] = append(dist.servers[i], servers[key])
		}
	}

	return
}
