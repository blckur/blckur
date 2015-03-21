package queue

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	clst *cluster
)

func update() {
	for {
		db := database.GetDatabase()
		coll := db.Nodes()
		nodes := []*nodes.Node{}
		cst := &cluster{
			servers: []string{},
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
			cst.servers = append(cst.servers, node.Id)
		}

		clst = cst
	}
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("beanstalkd", "update", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("settings", "queue", func(_ *messenger.Message) {
		go update()
	})
	update()

	requires.Register("queue")
}
