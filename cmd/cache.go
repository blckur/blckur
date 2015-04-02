package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts clustered redis cache node
func Cache() {
	opts := getServiceOptions()

	node := nodes.CacheNode{
		Id:   opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
