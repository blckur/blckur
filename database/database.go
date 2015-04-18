// Interface to mongodb database.
package database

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo"
	"os"
	"strings"
	"time"
)

var (
	MongoUrl string
	Session  *mgo.Session
)

type Database struct {
	session  *mgo.Session
	database *mgo.Database
}

func (d *Database) Close() {
	d.session.Close()
}

func (d *Database) getCollection(name string) (coll *Collection) {
	coll = &Collection{
		*d.database.C(name),
		d,
	}
	return
}

func (d *Database) Accounts() (coll *Collection) {
	coll = d.getCollection("accounts")
	return
}

func (d *Database) Notifications() (coll *Collection) {
	coll = d.getCollection("notifications")
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

func (d *Database) Tasks() (coll *Collection) {
	coll = d.getCollection("tasks")
	return
}

func (d *Database) Streams() (coll *Collection) {
	coll = d.getCollection("streams")
	return
}

func (d *Database) Nodes() (coll *Collection) {
	coll = d.getCollection("nodes")
	return
}

func (d *Database) Settings() (coll *Collection) {
	coll = d.getCollection("settings")
	return
}

func (d *Database) SessionKeys() (coll *Collection) {
	coll = d.getCollection("session_keys")
	return
}

func Connect() (err error) {
	Session, err = mgo.Dial(MongoUrl)
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

	var dbName string
	if x := strings.LastIndex(MongoUrl, "/"); x != -1 {
		dbName = MongoUrl[x+1:]
	} else {
		dbName = "blckur"
	}

	database := session.DB(dbName)

	db = &Database{
		session:  session,
		database: database,
	}
	return
}

func addIndexes() (err error) {
	db := GetDatabase()
	defer db.Close()

	coll := db.Accounts()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"user_id"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Notifications()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{
			"account_id",
			"remote_id",
		},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Notifications()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{
			"user_id",
			"timestamp",
		},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Users()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Messages()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"channel"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Tasks()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 6 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Streams()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 1 * time.Minute,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Nodes()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"type"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Nodes()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"timestamp"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.SessionKeys()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 72 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	return
}

func addCollections() (err error) {
	db := GetDatabase()
	defer db.Close()
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
		Capped:   true,
		MaxDocs:  1024,
		MaxBytes: 100000,
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func init() {
	MongoUrl = os.Getenv("DB")

	module := requires.New("database")

	module.Handler = func() {
		for {
			err := Connect()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Connection")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}

		for {
			err := addCollections()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Add collections")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}

		for {
			err := addIndexes()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Add indexes")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}
	}
}
