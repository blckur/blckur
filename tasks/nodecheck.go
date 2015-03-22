package tasks

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/nodes"
	"github.com/dropbox/godropbox/container/set"
	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
	"time"
)

type nodeCheck struct {
}

func (n *nodeCheck) Type() string {
	return "nodes"
}

func (n *nodeCheck) Run(db *database.Database) {
	coll := db.Nodes()

	nds := []nodes.Node{}
	err := coll.Find(bson.M{
		"timestamp": bson.M{
			"$lt": time.Now().Add(-5 * time.Minute),
		},
	}).All(&nds)
	if err != nil {
		err = database.ParseError(err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("tasks.nodecheck: Distinct error")
	}

	updated := set.NewSet()

	for _, node := range nds {
		err := coll.RemoveId(node.Id)
		if err != nil {
			err = database.ParseError(err)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("tasks.nodecheck: Remove error")
		}

		updated.Add(node.Type)
	}

	for nodeType := range updated.Iter() {
		messenger.Publish(db, nodeType.(string), "update")
	}
}

func init() {
	register(-1, -1, &nodeCheck{})
}
