package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/nodes"
)

func App() {
	logger.Init()
	database.Init()
	settings.Init()
	account.Init()
	cache.Init()
	messenger.Init()
	session.Init()
	opts := GetServiceOptions()

	node := nodes.AppNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
