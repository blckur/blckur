package streams

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Stream struct {
	Id bson.ObjectId `bson:"_id"`
	RunnerId bson.ObjectId `bson:"runner_id"`
	Type string `bson:"type"`
	Timestamp time.Time `bson:"timestamp"`
}
