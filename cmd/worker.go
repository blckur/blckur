package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/messenger"
	"github.com/Sirupsen/logrus"
	"time"
)

func Worker() {
	logger.Init()
	database.Init()
	settings.Init()
	queue.Init()
	messenger.Init()

	queue.NewListener(func (stream *queue.Stream) {
		for {
			job := stream.Reserve(30 * time.Second)
			if job == nil {
				return
			}

			err := job.Delete()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("worker: Job delete error")
			}
		}
	})

	block := make(chan int)
	<-block
}
