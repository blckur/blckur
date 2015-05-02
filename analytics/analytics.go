package analytics

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/search"
	"time"
)

type analyticEntry struct {
	Client    string `json:"client"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

func update() {
	conn := search.NewSession()
	if conn == nil {
		return
	}

	mapping := search.Mapping{}
	mapping.Add("client", true, "string", "not_analyzed")
	mapping.Add("timestamp", true, "date", "")
	mapping.Add("path", true, "string", "not_analyzed")

	err := conn.CreateIndex("analytics")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("analytics: Update error")
		time.Sleep(constants.RetryDelay)
		update()
		return
	}

	err = conn.PutMapping("analytics", "entry", analyticEntry{}, mapping)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("analytics: Update error")
		time.Sleep(constants.RetryDelay)
		update()
		return
	}
}

func init() {
	module := requires.New("logger.search")
	module.After("settings")

	module.Handler = func() {
		messenger.Register("settings", "search", func(_ *messenger.Message) {
			go update()
		})
		update()
	}
}
