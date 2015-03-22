package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/queue"
	"time"
)

func Worker() {
	utils.SeedRand()

	logger.Init()
	database.Init()
	settings.Init()
	queue.Init()

	queue.NewListener(func (stream *queue.Stream) {
		for {
			job := stream.Reserve(30 * time.Second)
			if job == nil {
				return
			}

			job.Delete()
		}
	})
}
