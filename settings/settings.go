package settings

import (
	"github.com/blckur/blckur/database"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"github.com/dropbox/godropbox/container/set"
	"github.com/blckur/blckur/utils"
)

var (
	System *system
	Twitter *twitter
)

type twitter struct {
	Id string `bson:"_id"`
	ConsumerKey string `bson:"consumer_key"`
	ConsumerSecret string  `bson:"consumer_secret"`
}

type system struct {
	Id string `bson:"_id"`
	CookieKey []byte `bson:"cookie_key"`
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

func Init() {
	utils.After("database")

	db := database.GetDatabase()
	defer db.Close()
	coll := db.Settings()

	System = &system{
		Id: "system",
	}
	err := parseFindError(coll.FindOneId("system", System))
	if err != nil {
		panic(err)
	}

	Twitter = &twitter{
		Id: "twitter",
	}
	err = parseFindError(coll.FindOneId("twitter", System))
	if err != nil {
		panic(err)
	}

	utils.Register("settings")
}
