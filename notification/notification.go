// Interface to notifications collection in database.
package notification

import (
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/database"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Notification struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId      bson.ObjectId `bson:"user_id" json:"-"`
	AccountId   bson.ObjectId `bson:"account_id" json:"-"`
	AccountType string        `bson:"account_type" json:"account_type"`
	RemoteId    string        `bson:"remote_id" json:"-"`
	Timestamp   time.Time     `bson:"timestamp" json:"timestamp"`
	Read        bool          `bson:"read" json:"read"`
	Type        string        `bson:"type,omitempty" json:"type"`
	Resource    string        `bson:"resource,omitempty" json:"-"`
	Origin      string        `bson:"origin,omitempty" json:"origin"`
	Link        string        `bson:"link,omitempty" json:"link"`
	Subject     string        `bson:"subject,omitempty" json:"subject"`
	Body        string        `bson:"body,omitempty" json:"body"`
}

func (n *Notification) Initialize(db *database.Database) (err error) {
	coll := db.Notifications()
	notf := &Notification{}

	coll.Find(bson.M{
		"user_id":    n.UserId,
		"account_id": n.AccountId,
		"remote_id":  n.RemoteId,
	}).Apply(mgo.Change{
		Update:    n,
		Upsert:    true,
		ReturnNew: true,
	}, notf)

	n.Id = notf.Id

	return
}

func (n *Notification) Commit(db *database.Database) (err error) {
	coll := db.Notifications()
	err = coll.Commit(n.Id, n)
	return
}

func (n *Notification) CommitFields(db *database.Database,
	fields set.Set) (err error) {

	coll := db.Notifications()
	err = coll.CommitFields(n.Id, n, fields)

	return
}

func GetNotification(db *database.Database, userId bson.ObjectId,
	notfId bson.ObjectId) (notf *Notification, err error) {

	coll := db.Notifications()
	notf = &Notification{}

	err = coll.FindOne(bson.M{
		"_id":     notfId,
		"user_id": userId,
	}, notf)

	return
}

func GetNotifications(db *database.Database, userId bson.ObjectId) (
	notfs []*Notification, err error) {

	coll := db.Notifications()
	notfs = []*Notification{}

	iter := coll.Find(bson.M{
		"user_id": userId,
		"type":    bson.M{"$exists": true},
	}).Limit(100).Sort("-timestamp").Iter()

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

func GetLastNotification(db *database.Database, userId bson.ObjectId,
	acctId bson.ObjectId) (notf *Notification, err error) {

	coll := db.Notifications()
	notf = &Notification{}

	err = coll.Find(bson.M{
		"user_id":    userId,
		"account_id": acctId,
	}).Sort("-timestamp").One(notf)
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			notf = nil
			err = nil
			return
		}
		return
	}
	return
}

type Publisher struct {
	UserIdHex string
	conn      *cache.ClusterConn
}

func (p *Publisher) New(data interface{}) (err error) {
	err = p.conn.Publish(p.UserIdHex, "notf_new", data)
	return
}

func (p *Publisher) Update(data interface{}) (err error) {
	err = p.conn.Publish(p.UserIdHex, "notf_update", data)
	return
}

func (p *Publisher) Close() {
	p.conn.Close()
}

func NewPublisher(userIdHex string) (pub *Publisher) {
	pub = &Publisher{
		UserIdHex: userIdHex,
		conn:      cache.Get(),
	}
	return
}
