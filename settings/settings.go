package settings

import (
    "github.com/blckur/blckur/database"
    "labix.org/v2/mgo/bson"
)

func Get(db *database.Database, group string, key string) (
        val interface{}, err error) {
    setCol := db.Settings()

    grp := map[string]interface{}{}

    err = setCol.Find(bson.M{
        "_id": group,
    }).Select(bson.M{
        key: 1,
    }).One(grp)
    if err != nil {
        return
    }

    val = grp[key]
    return
}

func Set(db *database.Database, group string, key string, val interface{}) (
        err error) {
    setCol := db.Settings()

    _, err = setCol.Upsert(bson.M{
        "_id": group,
    }, bson.M{
        key: val,
    })

    return
}
