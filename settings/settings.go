package settings

import (
    "github.com/blckur/blckur/database"
    "github.com/dropbox/godropbox/errors"
    "labix.org/v2/mgo/bson"
    "labix.org/v2/mgo"
)

func Get(db *database.Database, group string, key string) (
        val interface{}, err error) {
    setCol := db.Settings()

    grp := map[string]interface{}{}

    err = setCol.Find(&bson.M{
        "_id": group,
    }).Select(&bson.M{
        key: 1,
    }).One(grp)
    if err != nil {
        if err == mgo.ErrNotFound {
            err = nil
            return
        }

        err = &DatabaseError{
            errors.Wrap(err, "settings: Database error"),
        }
        return
    }

    val = grp[key]
    return
}

func Set(db *database.Database, group string, key string, val interface{}) (
        err error) {
    setCol := db.Settings()

    _, err = setCol.Upsert(&bson.M{
        "_id": group,
    }, &bson.M{
        key: val,
    })
    if err != nil {
        err = &DatabaseError{
            errors.Wrap(err, "settings: Database error"),
        }
        return
    }

    return
}
