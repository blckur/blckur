package notification

import (
	"github.com/blckur/blckur/database"
	"labix.org/v2/mgo/bson"
	"time"
)

type Notification struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId bson.ObjectId `bson:"user_id" json:"-"`
	RemoteId string `bson:"remote_id" json:"-"`
	Timestamp time.Time `bson:"timestamp" json"timestamp"`
	AccountType string `bson:"account_type" json:"account_type"`
	Type string `bson:"type,omitempty" json:"type"`
	Label string `bson:"-" json:"label"`
	Subject string `bson:"subject,omitempty" json:"subject"`
	Body string `bson:"body,omitempty" json:"body"`
}

func (n *Notification) Initialize(db *database.Database) (err error) {
	coll := db.Notifications()

	_, err = coll.Upsert(bson.M{
		"user_id": n.UserId,
		"remote_id": n.RemoteId,
	}, n)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetNotifications(db* database.Database, userId bson.ObjectId) (
		notfs []*Notification, err error) {
	coll := db.Notifications()
	notfs = []*Notification{}

	iter := coll.Find(bson.M{
		"user_id": userId,
	}).Sort("-timestamp").Iter()

	notf := &Notification{}
	for iter.Next(notf) {
		notfs = append(notfs, notf)
		notf = &Notification{}
	}

	err = iter.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
