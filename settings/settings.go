// Auto updating settings stored on mongodb.
package settings

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"reflect"
	"strconv"
	"time"
)

var (
	registry = map[string]interface{}{}
)

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

func setDefaults(obj interface{}) {
	val := reflect.ValueOf(obj)
	elm := val.Elem()

	n := elm.NumField()
	for i := 0; i < n; i++ {
		fld := elm.Field(i)
		typ := elm.Type().Field(i)

		if typ.PkgPath != "" {
			continue
		}

		tag := typ.Tag.Get("default")
		if tag == "" {
			continue
		}

		switch fld.Kind() {
		case reflect.Bool:
			parVal, err := strconv.ParseBool(tag)
			if err != nil {
				panic(err)
			}
			fld.SetBool(parVal)
		case reflect.Int:
			parVal, err := strconv.Atoi(tag)
			if err != nil {
				panic(err)
			}
			fld.SetInt(int64(parVal))
		case reflect.String:
			fld.SetString(tag)
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
		err = database.ParseError(err)
		return
	}

	setDefaults(data)

	return
}

func Update(name string) {
	group, ok := registry[name]
	if !ok {
		return
	}

	for {
		err := update(name, group)
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

func register(name string, group interface{}) {
	registry[name] = group
}

func init() {
	module := requires.New("settings")
	module.After("database")

	module.Handler = func() {
		for name, _ := range registry {
			Update(name)
		}

		messenger.Register("settings", "all", func(msg *messenger.Message) {
			Update(msg.Data.(string))
		})
	}
}
