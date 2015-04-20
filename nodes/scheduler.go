package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/scheduler"
)

type SchedulerNode struct {
	Id string
}

func (w *SchedulerNode) Start() {
	logrus.WithFields(logrus.Fields{
		"id": w.Id,
	}).Info("nodes.scheduler: Starting scheduler node")

	scheduler.Run()
}
