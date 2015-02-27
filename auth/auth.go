package auth

import (
    "github.com/blckur/blckur/user"
    "github.com/blckur/blckur/database"
    "labix.org/v2/mgo/bson"
    "time"
)

type Session struct {
    Id bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
    UserId bson.ObjectId `bson:"user_id" json:"user_id" binding:"required"`
    Timestamp time.Time `bson:"timestamp" json:"-"`
    db *database.Database
}

func (s *Session) Remove() (err error) {
    err = RemoveSession(s.db, s.Id)
    return
}

func (s *Session) GetUser() (usr *user.User, err error) {
    usr, err = user.GetUser(s.db, s.UserId)
    return
}

func GetSession(db *database.Database, id bson.ObjectId) (
        sess *Session, err error) {
    sessCol := db.Sessions()
    sess = &Session{}

    err = sessCol.FindId(id).One(sess)
    if err != nil {
        err = database.ParseError(err)
        return
    }

    sess.db = db

    return
}

func NewSession(db *database.Database, userId bson.ObjectId) (
        sess *Session, err error) {
    sessCol := db.Sessions()

    sess = &Session{
        Id: bson.NewObjectId(),
        UserId: userId,
        Timestamp: time.Now(),
        db: db,
    }

    err = sessCol.Insert(sess)
    if err != nil {
        err = database.ParseError(err)
        return
    }

    return
}

func RemoveSession(db *database.Database, id bson.ObjectId) (err error) {
    sessCol := db.Sessions()

    err = sessCol.RemoveId(id)
    if err != nil {
        err = database.ParseError(err)
        return
    }

    return
}
