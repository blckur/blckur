package auth

import (
    "github.com/blckur/blckur/database"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "github.com/dropbox/godropbox/errors"
    "time"
)

type Session struct {
    Id bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
    UserId bson.ObjectId `bson:"user_id" json:"user_id" binding:"required"`
    Timestamp time.Time `bson:"timestamp" json:"-"`
    db *database.Database
}

func (s *Session) Clear() (err error) {
    sessCol := s.db.Sessions()

    err = sessCol.RemoveId(s.Id)
    if err != nil {
        err = database.ParseError(err)
        return
    }

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
