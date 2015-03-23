package cmd

import (
	"github.com/blckur/blckur/scheduler"
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/nodes"
)

func Scheduler() {
	logger.Init()
	database.Init()
	settings.Init()
	scheduler.Init()
	queue.Init()
	messenger.Init()
	opts := GetServiceOptions()

	node := nodes.SchedulerNode{
		Id: opts.Id,
	}

	node.Start()
}
