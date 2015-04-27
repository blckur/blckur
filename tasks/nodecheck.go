package tasks

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/node"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"time"
)

type nodeCheck struct {
}

func (n *nodeCheck) Type() string {
	return "nodes"
}

func (n *nodeCheck) Run(db *database.Database) (err error) {
	coll := db.Nodes()

	nds := []node.Node{}
	err = coll.Find(bson.M{
		"timestamp": bson.M{
			"$lt": time.Now().Add(-5 * time.Minute),
		},
	}).All(&nds)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	updated := set.NewSet()

	for _, node := range nds {
		err = coll.RemoveId(node.Id)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		updated.Add(node.Type)
	}

	for nodeType := range updated.Iter() {
		messenger.Publish(db, nodeType.(string), "update")
	}

	return
}

func init() {
	register(-1, -1, 0, &nodeCheck{})
}
