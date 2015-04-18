package cmd

import (
	"github.com/blckur/blckur/nodes"
	"os"
)

// Starts web server node with JSON api
func App() {
	opts := getServiceOptions()

	source := os.Getenv("SOURCE")
	if source == "" {
		panic("Missing source")
	}

	node := nodes.AppNode{
		Id:         opts.Id,
		Host:       opts.Host,
		Address:    opts.Address,
		Port:       opts.Port,
		PublicPort: opts.PublicPort,
		Source:     source,
	}

	node.Start()
}
