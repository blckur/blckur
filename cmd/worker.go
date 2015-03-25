package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func Worker() {
	opts := getServiceOptions()

	node := nodes.WorkerNode{
		Id: opts.Id,
	}

	node.Start()
}
