package user

import (
    "github.com/blckur/blckur/utils"
    "github.com/blckur/blckur/database"
    "labix.org/v2/mgo/bson"
    "crypto/sha512"
    "math/rand"
    "time"
    "bytes"
)

type User struct {
    Id bson.ObjectId `bson:"_id,omitempty"`
    Email string `bson:"email"`
    PasswordSalt []byte `bson:"password_salt"`
    PasswordHash []byte `bson:"password_hash"`
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
    if (err != nil) {
        return
    }

    u.PasswordSalt = salt
    u.PasswordHash = u.hashPassword(password)

    return
}

func FindUser(email string) (usr *User, err error) {
    db := database.GetDatabase()
    usrs := db.Users()
    usr = &User{}

    err = usrs.Find(bson.M{
        "email": email,
    }).One(usr)

    return
}

func NewUser(email string, password string) (acct *User, err error) {
    db := database.GetDatabase()
    usrs := db.Users()

    usr := &User{
        Email: email,
    }
    usr.SetPassword(password)

    err = usrs.Insert(usr)

    return
}
