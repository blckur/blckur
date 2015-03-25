package nodes

import (
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
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

	db := database.GetDatabase()

	queue.NewListener(func (stream *queue.Stream) {
		for {
			job := stream.Reserve(30 * time.Second)
			if job == nil {
				return
			}

			if job.Type == "sync" {
				acct, err := account.GetAccount(db, "", job.Resource)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("worker: Sync job error")
				}

				client := acct.GetClient()
				err = client.Sync(db)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("worker: Sync job error")
				}
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
