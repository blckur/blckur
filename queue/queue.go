package queue

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/node"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/gdefer"
	"github.com/dropbox/godropbox/container/set"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"time"
	"sync"
)

var (
	clst *cluster
	listeners = []*Listener{}
	mutex = sync.Mutex{}
)

func Put(data interface{}, priority int,
		delay time.Duration, ttr time.Duration) (err error) {
	err = clst.Put(data, priority, delay, ttr)
	return
}

func update() {
	for {
		db := database.GetDatabase()
		coll := db.Nodes()
		nodes := []*node.Node{}

		err := coll.Find(bson.M{
			"type": "queue",
		}).All(&nodes)
		if err != nil {
			err = database.ParseError(err)
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		mutex.Lock()
		servers := set.NewSet()
		serversSlc := []string{}

		for _, node := range nodes {
			servers.Add(node.Address)
		}

		for server := range servers.Iter() {
			serversSlc = append(serversSlc, server.(string))
		}

		clst.servers = servers
		clst.serversSlc = serversSlc

		for _, lstnr := range listeners {
			lstnr.updateStreams(clst.servers)
		}
		mutex.Unlock()

		break
	}
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	clst = &cluster{
		defaultConsistency: settings.Queue.Consistency,
		servers: set.NewSet(),
		serversSlc: []string{},
		pool: map[string]*beanstalk.Conn{},
	}
	update()

	messenger.Register("queue", "update", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("settings", "queue", func(_ *messenger.Message) {
		clst.defaultConsistency = settings.Queue.Consistency
	})

	requires.Register("queue")

	gdefer.Defer(func() {
		clst.Close()
	})
}
