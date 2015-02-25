package auth

import (
    "labix.org/v2/mgo/bson"
    "github.com/blckur/blckur/database"
    "time"
)

type Session struct {
    Id bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
    UserId bson.ObjectId `bson:"user_id" json:"user_id" binding:"required"`
    Timestamp time.Time `bson:"timestamp" json:"-"`
}

func (s *Session) Clear() (err error) {
    db := database.GetDatabase()
    sessCol := db.Sessions()

    err = sessCol.RemoveId(s.Id)

    return
}

func NewSession(userId bson.ObjectId) (sess *Session, err error) {
    db := database.GetDatabase()
    sessCol := db.Sessions()

    sess = &Session{
        Id: bson.NewObjectId(),
        UserId: userId,
        Timestamp: time.Now(),
    }

    err = sessCol.Insert(sess)

    return
}
