package user

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/utils"
	"labix.org/v2/mgo/bson"
	"time"
)

type SessionKey struct {
	Id        string        `bson:"_id"`
	UserId    bson.ObjectId `bson:"user_id"`
	Timestamp time.Time     `bson:"timestamp"`
}

func (u *User) ResetPassword(db *database.Database) (
	key *SessionKey, err error) {

	coll := db.SessionKeys()
	key = &SessionKey{
		Id:        utils.RandStr(32),
		UserId:    u.Id,
		Timestamp: time.Now(),
	}

	err = coll.Insert(key)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	return
}

func GetSessionUser(db *database.Database, key string) (usr *User, err error) {
	coll := db.Users()
	sessColl := db.SessionKeys()

	sess := &SessionKey{}
	err = sessColl.FindOneId(key, sess)
	if err != nil {
		return
	}

	usr = &User{}
	err = coll.FindOneId(sess.UserId, usr)
	if err != nil {
		return
	}

	err = sessColl.RemoveId(sess.Id)
	if err != nil {
		return
	}

	return
}
