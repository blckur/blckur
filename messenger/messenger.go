// Pub/sub messaging system using mongodb tailable cursor.
package messenger

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

var (
	listeners = map[string][]func(*Message){}
)

type Message struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Channel   string        `bson:"channel"`
	Timestamp time.Time     `bson:"timestamp"`
	Data      interface{}   `bson:"data"`
}

func getCursorId(coll *database.Collection, channels []string) (
	id bson.ObjectId, err error) {

	msg := &Message{}

	var query *bson.M
	if len(channels) == 1 {
		query = &bson.M{
			"channel": channels[0],
		}
	} else {
		query = &bson.M{
			"channel": &bson.M{
				"$in": channels,
			},
		}
	}

	for i := 0; i < 2; i++ {
		err = coll.Find(query).Sort("-$natural").One(msg)

		if err != nil {
			err = database.ParseError(err)
			if i > 0 {
				return
			}

			switch err.(type) {
			case *database.NotFoundError:
				// Cannot use client-side ObjectId for tailable collection
				err = Publish(coll.Database, channels[0], nil)
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
		Channel:   channel,
		Timestamp: time.Now(),
		Data:      data,
	}

	err = coll.Insert(msg)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Subscribe(db *database.Database, channels []string,
	duration time.Duration, onMsg func(*Message) bool) (err error) {

	coll := db.Messages()
	cursorId, err := getCursorId(coll, channels)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	var channelBson interface{}
	if len(channels) == 1 {
		channelBson = channels[0]
	} else {
		channelBson = &bson.M{
			"$in": channels,
		}
	}

	query := &bson.M{
		"_id": &bson.M{
			"$gt": cursorId,
		},
		"channel": channelBson,
	}
	iter := coll.Find(query).Sort("$natural").Tail(duration)
	defer func() {
		iter.Close()
	}()

	for {
		msg := &Message{}
		for iter.Next(msg) {
			cursorId = msg.Id

			if msg.Data == nil {
				// Blank msg for cursor
				continue
			}

			if !onMsg(msg) {
				return
			}
		}

		if iter.Err() != nil {
			err = iter.Close()
			return
		}

		if iter.Timeout() {
			if !onMsg(nil) {
				return
			}
			continue
		}

		query := &bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channelBson,
		}
		iter = coll.Find(query).Sort("$natural").Tail(duration)
	}
}

func Register(channel string, event string, callback func(*Message)) {
	key := channel + ":" + event

	callbacks := listeners[key]

	if callbacks == nil {
		callbacks = []func(*Message){}
	}

	listeners[key] = append(callbacks, callback)
}

func init() {
	module := requires.New("messenger")
	module.After("settings")

	module.Handler = func() {
		go func() {
			channelsSet := set.NewSet()

			for key, _ := range listeners {
				channelsSet.Add(strings.Split(key, ":")[0])
			}

			channels := []string{}

			for channel := range channelsSet.Iter() {
				channels = append(channels, channel.(string))
			}

			for {
				db := database.GetDatabase()

				err := Subscribe(db, channels, 10*time.Second,
					func(msg *Message) bool {
						if msg == nil {
							return true
						}

						key := msg.Channel + ":all"
						for _, listener := range listeners[key] {
							listener(msg)
						}

						key = msg.Channel + ":" + msg.Data.(string)
						for _, listener := range listeners[key] {
							listener(msg)
						}

						return true
					})
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("messenger: Listener")
				}

				time.Sleep(constants.RETRY_DELAY)
			}
		}()
	}
}
