package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/scheduler"
)

type SchedulerNode struct {
	Id string
}

func (w *SchedulerNode) Start() {
	constants.Node = "scheduler"

	logrus.WithFields(logrus.Fields{
		"id":      w.Id,
		"version": constants.Version,
	}).Info("nodes.scheduler: Starting scheduler node")

	scheduler.Run()
}
