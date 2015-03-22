package cmd

import (
	"github.com/blckur/blckur/scheduler"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
)

func Scheduler() {
	logger.Init()
	database.Init()
	settings.Init()
	scheduler.Init()
	messenger.Init()

	scheduler.Run()
}
