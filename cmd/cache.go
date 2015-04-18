package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts clustered redis cache node
func Cache() {
	opts := getServiceOptions()

	node := nodes.CacheNode{
		Id:         opts.Id,
		Host:       opts.Host,
		Address:    opts.Address,
		Port:       opts.Port,
		PublicPort: opts.PublicPort,
	}

	node.Start()
}
