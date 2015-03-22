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
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"time"
	"sync"
)

var (
	clst *cluster
	mutex sync.Mutex
	listeners []*Listener
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
		nodes := []*nodes.Node{}

		err := coll.Find(bson.M{
			"type": "beanstalkd",
		}).All(&nodes)
		if err != nil {
			err = database.ParseError(err)
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		mutex.Lock()
		for _, node := range nodes {
			clst.servers.Add(node.Address)
		}

		servers := []string{}
		for server := range clst.servers.Iter() {
			servers = append(servers, server.(string))
		}
		clst.serversSlc = servers

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

	messenger.Register("beanstalkd", "update", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("settings", "queue", func(_ *messenger.Message) {
		go update()
	})

	requires.Register("queue")

	gdefer.Defer(func() {
		clst.Close()
	})
}
