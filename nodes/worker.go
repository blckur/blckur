package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/utils"
	"time"
)

type WorkerNode struct {
	Id string
}

func (w *WorkerNode) sync(db *database.Database, job *queue.Job) (err error) {
	acct, err := account.GetAccount(db, "", job.Resource)
	if err != nil {
		return
	}

	client, err := acct.GetClient()
	if err != nil {
		return
	}

	err = client.Sync(db)
	if err != nil {
		return
	}

	return
}

func (w *WorkerNode) Start() {
	logrus.WithFields(logrus.Fields{
		"id":      w.Id,
		"version": utils.GetVersion(),
	}).Info("nodes.worker: Starting worker node")

	db := database.GetDatabase()

	queue.NewListener(func(stream *queue.Stream) {
		for {
			job := stream.Reserve(30 * time.Second)
			if job == nil {
				return
			}

			if job.Type == "sync" {
				err := w.sync(db, job)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("worker: Sync job error")
				}
			}

			job.Delete()
		}
	})

	block := make(chan int)
	<-block
}
