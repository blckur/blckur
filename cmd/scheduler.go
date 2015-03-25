package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts task scheduler node
func Scheduler() {
	opts := getServiceOptions()

	node := nodes.SchedulerNode{
		Id: opts.Id,
	}

	node.Start()
}
