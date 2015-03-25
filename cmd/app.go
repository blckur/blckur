package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func App() {
	opts := getServiceOptions()

	node := nodes.AppNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
