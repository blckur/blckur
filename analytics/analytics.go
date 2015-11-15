package analytics

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/search"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	buffer chan *Entry
)

type Entry struct {
	Client    string `json:"client"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

func (e *Entry) Send() {
	if len(buffer) <= 1000 {
		buffer <- e
	}
}

func listen() {
	for {
		entry := <-buffer

		conn := search.NewSession()
		if conn == nil {
			continue
		}

		err := conn.Index("analytics", "entry",
			bson.NewObjectId().Hex(), entry)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("analytics: Send error")
		}
	}
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

	err = conn.PutMapping("analytics", "entry", Entry{}, mapping)
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
	buffer = make(chan *Entry, 1024)

	module := requires.New("analytics")
	module.After("settings")

	module.Handler = func() {
		messenger.Register("settings", "search", func(_ *messenger.Message) {
			go update()
		})
		update()
	}

	go listen()
}
