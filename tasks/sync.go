package tasks

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/queue"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo/bson"
	"time"
)

type sync struct {
}

func (s *sync) Type() string {
	return "sync"
}

func (s *sync) Run(db *database.Database) (err error) {
	coll := db.Accounts()

	acctIds := []bson.ObjectId{}
	err = coll.Find(bson.M{}).Distinct("_id", &acctIds)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	for _, acctId := range acctIds {
		job := queue.NewJob(time.Duration(settings.Task.Timeout) * time.Second)
		job.Type = s.Type()
		job.Resource = acctId

		err = job.Put(0, 0,
			time.Duration(settings.Task.RetryTimeout) * time.Second)
		if err != nil {
			return
		}
	}

	return
}

func init() {
	register(-1, -1, &sync{})
}
