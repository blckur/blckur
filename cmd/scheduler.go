package cmd

import (
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/nodes"
)

func Scheduler() {
	requires.Init()
	opts := GetServiceOptions()

	node := nodes.SchedulerNode{
		Id: opts.Id,
	}

	node.Start()
}
