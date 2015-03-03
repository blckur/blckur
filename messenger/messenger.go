package messenger

import (
	"github.com/blckur/blckur/database"
	"labix.org/v2/mgo/bson"
	"time"
	"fmt"
)

type Message struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Channel string `bson:"channel"`
	Timestamp time.Time `bson:"timestamp"`
	Data interface{} `bson:"data"`
}

func getCursorId(coll *database.Collection, channel string) (
		id bson.ObjectId, err error) {
	msg := &Message{}

	for i := 0; i < 2; i++ {
		err = coll.Find(&bson.M{
			"channel": channel,
		}).Sort("-$natural").One(msg)

		if err != nil {
			err = database.ParseError(err)
			if i > 0 {
				return
			}

			switch err.(type) {
			case *database.NotFoundError:
				// Cannot use client-side ObjectId for tailable collection
				err = Publish(coll.Database, channel, nil)
				if err != nil {
					err = database.ParseError(err)
					return
				}
				continue
			default:
				return
			}
		} else {
			break
		}
	}

	id = msg.Id

	return
}

func Publish(db *database.Database, channel string, data interface{}) (
		err error) {
	coll := db.Messages()

	msg := &Message{
		Channel: channel,
		Timestamp: time.Now(),
		Data: data,
	}

	err = coll.Insert(msg)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Subscribe(db *database.Database, channel string, duration time.Duration,
		onMsg func(*Message) bool) (err error) {
	coll := db.Messages()
	cursorId, err := getCursorId(coll, channel)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	iter := coll.Find(&bson.M{
		"_id": &bson.M{
			"$gt": cursorId,
		},
		"channel": channel,
	}).Sort("$natural").Tail(duration)
	defer func() {
		iter.Close()
	}()

	for {
		msg := &Message{}
		for iter.Next(msg) {
			if msg.Data == nil {
				// Blank msg for cursor
				continue
			}

			if onMsg(msg) == false {
				return
			}
		}

		if iter.Err() != nil {
			err = iter.Close()
			return
		}

		if iter.Timeout() {
			if onMsg(nil) == false {
				return
			}
			continue
		}

		iter = coll.Find(&bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channel,
		}).Sort("$natural").Tail(duration)
	}

	return
}
