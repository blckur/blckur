package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/cache"
)

func Worker() {
	logger.Init()
	database.Init()
	settings.Init()
	queue.Init()
	cache.Init()
	messenger.Init()
	opts := GetServiceOptions()

	node := nodes.WorkerNode{
		Id: opts.Id,
	}

	node.Start()
}
