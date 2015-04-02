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
		Id:     opts.Id,
		Host:   opts.Host,
		Port:   opts.Port,
		Source: source,
	}

	node.Start()
}
