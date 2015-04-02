// Stores sessions in cookies.
package session

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/user"
	"labix.org/v2/mgo/bson"
	"time"
)

type Session struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId    bson.ObjectId `bson:"user_id" json:"user_id"`
	Timestamp time.Time     `bson:"timestamp" json:"-"`
	coll      *database.Collection
}

func (s *Session) Remove() (err error) {
	err = RemoveSession(s.coll.Database, s.Id)
	return
}

func (s *Session) GetUser() (usr *user.User, err error) {
	usr, err = user.GetUser(s.coll.Database, s.UserId)
	return
}

func GetSession(db *database.Database, id bson.ObjectId) (
	sess *Session, err error) {
	coll := db.Sessions()
	sess = &Session{}

	err = coll.FindOneId(id, sess)
	if err != nil {
		return
	}

	sess.coll = coll

	return
}

func NewSession(db *database.Database, userId bson.ObjectId) (
	sess *Session, err error) {
	coll := db.Sessions()
	sess = &Session{
		Id:        bson.NewObjectId(),
		UserId:    userId,
		Timestamp: time.Now(),
		coll:      coll,
	}

	err = coll.Insert(sess)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveSession(db *database.Database, id bson.ObjectId) (err error) {
	coll := db.Sessions()

	err = coll.RemoveId(id)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
