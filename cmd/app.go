package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func App() {
	opts := GetServiceOptions()

	node := nodes.AppNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
