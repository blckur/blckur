package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts worker node listening to queue
func Worker() {
	opts := getServiceOptions()

	node := nodes.WorkerNode{
		Id: opts.Id,
	}

	node.Start()
}
