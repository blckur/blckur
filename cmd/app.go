package cmd

import (
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/requires"
)

func App() {
	requires.Init()


	opts := GetServiceOptions()

	node := nodes.AppNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
