// Beanstalkd cluster with client side sharding/replication.
package queue

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/gdefer"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/node"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/container/set"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

var (
	clst      *cluster
	listeners = []*Listener{}
	mutex     = sync.Mutex{}
)

func put(job *Job, priority int, delay time.Duration,
	ttr time.Duration) (err error) {
	err = clst.Put(job, priority, delay, ttr)
	return
}

func update() {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Nodes()
	nodes := []*node.Node{}

	err := coll.Find(bson.M{
		"type": "queue",
	}).All(&nodes)
	if err != nil {
		err = database.ParseError(err)
		time.Sleep(constants.RetryDelay)
		update()
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
}

func init() {
	module := requires.New("queue")
	module.After("settings")
	module.Before("messenger")

	module.Handler = func() {
		clst = &cluster{
			defaultConsistency: settings.Queue.Consistency,
			servers:            set.NewSet(),
			serversSlc:         []string{},
			pool:               map[string]*beanstalk.Conn{},
		}
		update()

		messenger.Register("queue", "update", func(_ *messenger.Message) {
			go update()
		})
		messenger.Register("settings", "queue", func(_ *messenger.Message) {
			clst.defaultConsistency = settings.Queue.Consistency
		})

		gdefer.Defer(func() {
			clst.Close()
		})
	}
}
