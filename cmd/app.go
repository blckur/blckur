package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts web server node with JSON api
func App() {
	opts := getServiceOptions()

	node := nodes.AppNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
