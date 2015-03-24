package cmd

import (
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/nodes"
)

func Worker() {
	requires.Init()
	opts := GetServiceOptions()

	node := nodes.WorkerNode{
		Id: opts.Id,
	}

	node.Start()
}
