package user

import (
	"bytes"
	"crypto/sha512"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
	Email        string        `bson:"email" json:"email" binding:"required"`
	PasswordSalt []byte        `bson:"password_salt" json:"-"`
	PasswordHash []byte        `bson:"password_hash" json:"-"`
	db           *database.Database
}

func (u *User) hashPassword(password string) []byte {
	hash := sha512.New()
	hash.Write(u.PasswordSalt)
	hash.Write([]byte(password))
	digest := hash.Sum(nil)

	for i := 0; i < 5; i++ {
		hash = sha512.New()
		hash.Write(digest)
		digest = hash.Sum(nil)
	}

	return digest
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
	usrsCol := db.Users()

	err = usrsCol.UpdateId(u.Id, bson.M{
		"$set": u,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func FindUser(db *database.Database, email string) (usr *User, err error) {
	usrsCol := db.Users()
	usr = &User{
		db: db,
	}

	err = usrsCol.Find(bson.M{
		"email": email,
	}).One(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetUser(db *database.Database, id bson.ObjectId) (usr *User, err error) {
	usrsCol := db.Users()
	usr = &User{
		db: db,
	}

	err = usrsCol.FindId(id).One(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewUser(db *database.Database, email string, password string) (
	usr *User, err error) {
	usrsCol := db.Users()

	usr = &User{
		Id:    bson.NewObjectId(),
		Email: email,
		db:    db,
	}
	usr.SetPassword(password)

	err = usrsCol.Insert(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
