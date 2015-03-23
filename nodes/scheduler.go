package nodes

import (
	"github.com/blckur/blckur/scheduler"
	"github.com/Sirupsen/logrus"
)

type SchedulerNode struct {
	Id string
}

func (w *SchedulerNode) Start() {
	logrus.WithFields(logrus.Fields{
		"id": w.Id,
	}).Info("nodes.worker: Starting scheduler node")

	scheduler.Run()
}
