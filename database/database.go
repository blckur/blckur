package database

import (
    "labix.org/v2/mgo"
    "github.com/dropbox/godropbox/errors"
)

var Session *mgo.Session

type Database struct {
    database *mgo.Database
}

type Collection struct {
    *mgo.Collection
}

func (d *Database) Users() (coll *Collection) {
    coll = &Collection{
        d.database.C("users"),
    }
    return
}

func (d *Database) Sessions() (coll *Collection) {
    coll = &Collection{
        d.database.C("sessions"),
    }
    return
}

func (d *Database) Settings() (coll *Collection) {
    coll = &Collection{
        d.database.C("settings"),
    }
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
        database: database,
    }
    return
}

func AddIndexes() (err error) {
    db := GetDatabase()

    usrsCol := db.Users()
    err = usrsCol.EnsureIndex(mgo.Index{
        Key: []string{"email"},
        Unique: true,
        Background: true,
    })
    if err != nil {
        err = &IndexError{
            errors.Wrap(err, "database: Index error"),
        }
    }

    return
}

func Init() (err error) {
    err = Connect()
    if err != nil {
        return
    }

    err = AddIndexes()

    return
}

func GetErrorCode(err error) (errCode int) {
    switch err := err.(type) {
    case *mgo.LastError:
        errCode = err.Code;
    case *mgo.QueryError:
        errCode = err.Code;
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
            errors.New("database: Unknown error"),
        }
    }

    return
}
