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

func (s *Stream) Update() (stop bool, err error) {
	coll := s.db.Streams()
	stop = false

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"timestamp": time.Now(),
			},
		},
		ReturnNew: true,
	}
	doc := &Stream{}

	_, err = coll.Find(bson.M{
		"_id": s.Id,
		"runner_id": s.RunnerId,
	}).Apply(change, doc)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			stop = true
			err = nil
		default:
			return
		}
	}

	s.Timestamp = doc.Timestamp

	return
}

func (s *Stream) initialize() (err error) {
	coll := s.db.Streams()

	s.Timestamp = time.Now()
	s.RunnerId = bson.NewObjectId()

	_, err = coll.UpsertId(s.Id, s)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
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
