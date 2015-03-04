package settings

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
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
	PapperTrail = &papperTrail{
		Id: "paper_trail",
	}
	Twitter = &twitter{
		Id: "twitter",
	}
	System = &system{
		Id: "system",
	}
}

type papperTrail struct {
	Id string `bson:"_id"`
	Address string `bson:"address"`
}

type twitter struct {
	Id string `bson:"_id"`
	ConsumerKey string `bson:"consumer_key"`
	ConsumerSecret string  `bson:"consumer_secret"`
}

type system struct {
	Id string `bson:"_id"`
	CookieKey []byte `bson:"cookie_key"`
	Domain string `bson:"domain"`
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

func Update(groupName string) {
	var group interface{}

	switch groupName {
	case "paper_trail":
		group = PapperTrail
	case "system":
		group = System
	case "twitter":
		group = Twitter
	}

	for {
		err := update(groupName, group)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("database: Update")
		} else {
			break
		}

		time.Sleep(constants.RETRY_DELAY)
	}
}

func Init() {
	utils.After("database")

	Update("paper_trail")
	Update("twitter")
	Update("system")

	messenger.Register("settings", "all", func(msg *messenger.Message) {
		Update(msg.Data.(string))
	})

	utils.Register("settings")
}
