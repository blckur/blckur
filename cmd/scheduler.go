package cmd

import (
	"github.com/blckur/blckur/nodes"
)

func Scheduler() {
	opts := getServiceOptions()

	node := nodes.SchedulerNode{
		Id: opts.Id,
	}

	node.Start()
}
