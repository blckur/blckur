package cmd

import (
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/requires"
)

func Queue() {
	requires.Init()
	opts := GetServiceOptions()

	node := nodes.QueueNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
