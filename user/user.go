// Interface to users collection in database.
package user

import (
	"bytes"
	"crypto/sha512"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Email string `bson:"email" json:"email"`
	PasswordSalt []byte `bson:"pass_salt" json:"-"`
	PasswordHash []byte `bson:"pass_hash" json:"-"`
}

func (u *User) hashPassword(password string) (digest []byte) {
	hash := sha512.New()
	hash.Write(u.PasswordSalt)
	hash.Write([]byte(password))
	digest = hash.Sum(nil)

	for i := 0; i < 5; i++ {
		hash = sha512.New()
		hash.Write(digest)
		digest = hash.Sum(nil)
	}

	return
}

func (u *User) CheckPassword(password string) bool {
	hash := u.hashPassword(password)
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Microsecond)
	return bytes.Equal(hash, u.PasswordHash)
}

func (u *User) SetPassword(password string) (err error) {
	salt, err := utils.RandBytes(16)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "user: Unknown error"),
		}
		return
	}

	u.PasswordSalt = salt
	u.PasswordHash = u.hashPassword(password)

	return
}

func (u *User) Commit(db *database.Database) (err error) {
	coll := db.Users()
	err = coll.Commit(u.Id, u)
	return
}

func (u *User) CommitFields(db *database.Database, fields set.Set) (
		err error) {
	coll := db.Users()
	err = coll.CommitFields(u.Id, u, fields)
	return
}

func FindUser(db *database.Database, email string) (usr *User, err error) {
	coll := db.Users()
	usr = &User{}

	err = coll.FindOne(bson.M{
		"email": email,
	}, usr)
	if err != nil {
		return
	}

	return
}

func GetUser(db *database.Database, id bson.ObjectId) (usr *User, err error) {
	coll := db.Users()
	usr = &User{}

	err = coll.FindOneId(id, usr)
	if err != nil {
		return
	}

	return
}

func NewUser(db *database.Database, email string, password string) (
		usr *User, err error) {
	coll := db.Users()
	usr = &User{
		Id: bson.NewObjectId(),
		Email: email,
		coll: coll,
	}

	usr.SetPassword(password)

	err = coll.Insert(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
