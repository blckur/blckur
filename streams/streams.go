package streams

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Stream struct {
	Id bson.ObjectId `bson:"_id"`
	Type string `bson:"type"`
	Timestamp time.Time `bson:"timestamp"`
}
