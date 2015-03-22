package cmd

import (
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/messenger"
)

func Queue() {
	logger.Init()
	database.Init()
	settings.Init()
	messenger.Init()
	opts := GetServiceOptions()

	node := nodes.QueueNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
