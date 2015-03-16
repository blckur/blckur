package streams

import (
	"github.com/blckur/blckur/database"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Stream struct {
	Id bson.ObjectId `bson:"_id"`
	RunnerId bson.ObjectId `bson:"runner_id"`
	Type string `bson:"type"`
	Timestamp time.Time `bson:"timestamp"`
	db *database.Database
}

func NewStream(acctId bson.ObjectId, typ string) (stream *Stream) {
	db := database.GetDatabase()

	stream = &Stream{
		Id: acctId,
		Type: typ,
		db: db,
	}

	return
}
