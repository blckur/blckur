package database

import (
    "labix.org/v2/mgo"
)

var Session *mgo.Session

type Database struct {
    database *mgo.Database
}

func (d *Database) Users() (coll *mgo.Collection) {
    coll = d.database.C("users")
    return
}

func Connect() (err error) {
    Session, err = mgo.Dial("localhost")
    if err != nil {
        return
    }
    Session.SetMode(mgo.Monotonic, true)
    return
}

func GetDatabase() (db *Database) {
    session := Session.Copy()
    database := session.DB("blckur")

    db = &Database{
        database: database,
    }
    return
}

func InitIndexes() (err error) {
    db := GetDatabase()

    usrsCol := db.Users()
    err = usrsCol.EnsureIndex(mgo.Index{
        Key: []string{"email"},
        Unique: true,
        Background: true,
    })
    if err != nil {
        return
    }

    return
}
