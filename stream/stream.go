// Data stream for streaming apis.
package stream

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Backend interface {
	Run()
	Stop()
}

type Stream struct {
	Id        bson.ObjectId `bson:"_id"`
	RunnerId  bson.ObjectId `bson:"runner_id"`
	Timestamp time.Time     `bson:"timestamp"`
	backend   Backend
	db        *database.Database
}

func (s *Stream) Start() (err error) {
	coll := s.db.Streams()

	count, err := coll.FindId(s.Id).Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if count > 0 {
		return
	}

	stop := false

	err = s.initialize()
	if err != nil {
		return
	}

	go func() {
		s.backend.Run()
		stop = true
	}()

	go func() {
		lastUpdate := time.Now()
		for {
			if stop {
				break
			}

			time.Sleep(time.Duration(
				settings.Stream.RefreshRate) * time.Second)
			stop, err := s.Update()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":      err,
					"account_id": s.Id.Hex(),
					"runner_id":  s.RunnerId.Hex(),
				}).Error("streams: Update error")
			} else if stop {
				s.backend.Stop()
				break
			} else {
				lastUpdate = time.Now()
			}

			if time.Since(lastUpdate) > 2*time.Minute {
				logrus.WithFields(logrus.Fields{
					"account_id": s.Id.Hex(),
					"runner_id":  s.RunnerId.Hex(),
				}).Error("streams: Update timed out")
				s.backend.Stop()
				break
			}
		}
	}()

	return
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
		"_id":       s.Id,
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

func NewStream(db *database.Database, acctId bson.ObjectId,
	backend Backend) (stream *Stream) {

	stream = &Stream{
		backend: backend,
		db:      db,
		Id:      acctId,
	}

	return
}
