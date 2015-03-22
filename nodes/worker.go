package nodes

import (
	"github.com/blckur/blckur/queue"
	"github.com/Sirupsen/logrus"
	"time"
)

type WorkerNode struct {
	Id string
}

func (w *WorkerNode) Start() {
	logrus.WithFields(logrus.Fields{
		"id": w.Id,
	}).Info("nodes.worker: Starting worker node")

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
