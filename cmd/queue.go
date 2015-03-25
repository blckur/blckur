package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func Queue() {
	opts := getServiceOptions()

	node := nodes.QueueNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
