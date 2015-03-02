package database

import (
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo"
	"fmt"
)

var Session *mgo.Session

type Database struct {
	session *mgo.Session
	database *mgo.Database
}

func (d *Database) Close() {
	d.session.Close()
}

func (d *Database) getCollection(name string) (coll *Collection) {
	coll = &Collection{
		d.database.C(name),
		d,
	}
	return
}

func (d *Database) Accounts() (coll *Collection) {
	coll = d.getCollection("accounts")
	return
}

func (d *Database) Tokens() (coll *Collection) {
	coll = d.getCollection("tokens")
	return
}

func (d *Database) Messages() (coll *Collection) {
	coll = d.getCollection("messages")
	return
}

func (d *Database) Users() (coll *Collection) {
	coll = d.getCollection("users")
	return
}

func (d *Database) Sessions() (coll *Collection) {
	coll = d.getCollection("sessions")
	return
}

func (d *Database) Settings() (coll *Collection) {
	coll = d.getCollection("settings")
	return
}

func Connect() (err error) {
	Session, err = mgo.Dial("localhost")
	if err != nil {
		err = &ConnectionError{
			errors.Wrap(err, "database: Connection error"),
		}
		return
	}

	Session.SetMode(mgo.Monotonic, true)

	return
}

func GetDatabase() (db *Database) {
	session := Session.Copy()
	database := session.DB("blckur")

	db = &Database{
		session:  session,
		database: database,
	}
	return
}

func AddIndexes() (err error) {
	db := GetDatabase()

	coll := db.Users()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{"email"},
		Unique: true,
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Messages()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{"channel"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	return
}

func AddCollections() (err error) {
	db := GetDatabase()
	coll := db.Messages()

	names, err := db.database.CollectionNames()
	if err != nil {
		err = ParseError(err)
		return
	}

	for _, name := range names {
		if name == "messages" {
			return
		}
	}

	err = coll.Create(&mgo.CollectionInfo{
		Capped: true,
		MaxDocs: 1024,
		MaxBytes: 100000,
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func Init() (err error) {
	err = Connect()
	if err != nil {
		return
	}

	err = AddCollections()
	if err != nil {
		return
	}

	err = AddIndexes()
	if err != nil {
		return
	}

	return
}

func GetErrorCode(err error) (errCode int) {
	switch err := err.(type) {
	case *mgo.LastError:
		errCode = err.Code
	case *mgo.QueryError:
		errCode = err.Code
	}

	return
}

func ParseError(err error) (newErr error) {
	if err == mgo.ErrNotFound {
		newErr = &NotFoundError{
			errors.New("database: Not found"),
		}
		return
	}

	errCode := GetErrorCode(err)

	switch errCode {
	case 11000, 11001, 12582, 16460:
		newErr = &DuplicateKeyError{
			errors.New("database: Duplicate key"),
		}
	default:
		newErr = &UnknownError{
			errors.New(fmt.Sprintf("database: Unknown error %d", errCode)),
		}
	}

	return
}
