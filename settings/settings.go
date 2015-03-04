package settings

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/utils"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	PapperTrail *papperTrail
	System *system
	Twitter *twitter
)

func init() {
	PapperTrail = &papperTrail{}
	Twitter = &twitter{}
	System = &system{}
}

type Group interface {
	Update() error
}

type papperTrail struct {
	Id string `bson:"_id"`
	Address string `bson:"address"`
}

func (t *papperTrail) Update() (err error) {
	err = update("paper_trail", t)
	return
}

type twitter struct {
	Id string `bson:"_id"`
	ConsumerKey string `bson:"consumer_key"`
	ConsumerSecret string  `bson:"consumer_secret"`
}

func (t *twitter) Update() (err error) {
	err = update("twitter", t)
	return
}

type system struct {
	Id string `bson:"_id"`
	CookieKey []byte `bson:"cookie_key"`
	Domain string `bson:"domain"`
}

func (t *system) Update() (err error) {
	err = update("system", t)
	return
}

func Commit(db *database.Database, group interface{}, fields set.Set) (
		err error) {
	coll := db.Settings()

	selector := database.SelectFields(group, set.NewSet("_id"))
	update := database.SelectFields(group, fields)

	_, err = coll.Upsert(selector, update)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Get(db *database.Database, group string, key string) (
		val interface{}, err error) {
	coll := db.Settings()

	grp := map[string]interface{}{}

	err = coll.Find(bson.M{
		"_id": group,
	}).Select(bson.M{
		key: 1,
	}).One(grp)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
			return
		default:
			err = &DatabaseError{
				errors.Wrap(err, "settings: Database error"),
			}
			return
		}
	}

	val = grp[key]
	return
}

func Set(db *database.Database, group string, key string, val interface{}) (
		err error) {
	coll := db.Settings()

	_, err = coll.Upsert(bson.M{
		"_id": group,
	}, bson.M{"$set": bson.M{
		key: val,
	}})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func parseFindError(inErr error) (err error) {
	if inErr != nil {
		switch inErr.(type) {
			case *database.NotFoundError:
			err = nil
			default:
			err = inErr
		}
	}

	return
}

func update(group string, data interface{}) (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Settings()

	err = parseFindError(coll.FindOneId(group, data))
	if err != nil {
		database.ParseError(err)
		return
	}

	return
}

func Init() {
	utils.After("database")

	System = &system{}
	for {
		err := System.Update()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("database: System update")
		} else {
			break
		}

		time.Sleep(constants.DB_RETRY_DELAY)
	}

	Twitter = &twitter{}
	for {
		err := Twitter.Update()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("database: Twitter update")
		} else {
			break
		}

		time.Sleep(constants.DB_RETRY_DELAY)
	}

	utils.Register("settings")
}
