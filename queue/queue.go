package queue

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/gdefer"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"time"
	"sync"
)

var (
	clst *cluster
)

func Put(data interface{}, priority int,
		delay time.Duration, ttr time.Duration) (err error) {
	err = clst.Put(data, priority, delay, ttr)
	return
}

func GetStreams() (streams []*Stream) {
	mutex.RLock()
	streams = clst.GetStreams()
	mutex.RUnlock()
	return
}

func update() {
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	for {
		db := database.GetDatabase()
		coll := db.Nodes()
		nodes := []*nodes.Node{}
		clst := &cluster{
			servers: set.NewSet(),
			defaultConsistency: settings.Queue.Consistency,
		}

		err := coll.Find(bson.M{
			"type": "queue",
		}).All(&nodes)
		if err != nil {
			err = database.ParseError(err)
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		for _, node := range nodes {
			clst.servers.Add(node.Id)
		}
	}

	messenger.Register("beanstalkd", "update", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("settings", "queue", func(_ *messenger.Message) {
		go update()
	})
	update()

	requires.Register("queue")

	gdefer.Defer(func() {
		clst.Close()
	})
}
