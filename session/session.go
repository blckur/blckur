package session

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/user"
	"labix.org/v2/mgo/bson"
	"time"
)

type Data struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
	UserId bson.ObjectId `bson:"user_id" json:"user_id" binding:"required"`
	Timestamp time.Time `bson:"timestamp" json:"-"`
}

type Session struct {
	*Data
	coll *database.Collection
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
	sess = &Session{
		&Data{},
		coll,
	}

	err = coll.FindOneId(id, sess)
	return
}

func NewSession(db *database.Database, userId bson.ObjectId) (
		sess *Session, err error) {
	coll := db.Sessions()
	sess = &Session{
		&Data{
			Id: bson.NewObjectId(),
			UserId: userId,
			Timestamp: time.Now(),
		},
		coll,
	}

	err = coll.Insert(sess.Data)
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
