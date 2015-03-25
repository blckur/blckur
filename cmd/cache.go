package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func Cache() {
	opts := getServiceOptions()

	node := nodes.CacheNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
