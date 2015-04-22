package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/scheduler"
	"github.com/blckur/blckur/utils"
)

type SchedulerNode struct {
	Id string
}

func (w *SchedulerNode) Start() {
	logrus.WithFields(logrus.Fields{
		"id":      w.Id,
		"version": utils.GetVersion(),
	}).Info("nodes.scheduler: Starting scheduler node")

	scheduler.Run()
}
