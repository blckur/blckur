package auth

import (
    "github.com/blckur/blckur/database"
    "github.com/blckur/blckur/settings"
    "github.com/blckur/blckur/errortypes"
    "github.com/gorilla/securecookie"
    "github.com/gorilla/sessions"
    "labix.org/v2/mgo/bson"
    "github.com/gin-gonic/gin"
    "github.com/dropbox/godropbox/errors"
)

var Store *sessions.CookieStore

type Cookie struct {
    Id bson.ObjectId
    data *sessions.Session
    con *gin.Context
}

func (c *Cookie) Get(key string) (val string) {
    valInter := c.data.Values[key]
    if valInter == nil {
        val = ""
    } else {
        val = valInter.(string)
    }
    return
}

func (c *Cookie) Set(key string, val string) {
    c.data.Values[key] = val
}

func (c *Cookie) GetSession(db *database.Database) (
        sess *Session, err error) {
    sessId := c.Get("id")
    if sessId == "" {
        err = &NotFoundError{
            errors.New("auth: Session not found"),
        }
        return
    }

    sess, err = GetSession(db, bson.ObjectIdHex(sessId))
    if err != nil {
        switch err.(type) {
        case *database.NotFoundError:
        default:
            err = &errortypes.UnknownError{
                errors.Wrap(err, "auth: Unknown session error"),
            }
        }
        return
    }

    return
}

func (c *Cookie) NewSession(db *database.Database, id bson.ObjectId) (
        sess *Session, err error) {
    sess, err = NewSession(db, id)
    if err != nil {
        err = &errortypes.UnknownError{
            errors.Wrap(err, "auth: Unknown session error"),
        }
        return
    }

    c.Set("id", sess.Id.Hex())
    err = c.Save()

    return
}

func (c *Cookie) Save() (err error) {
    err = c.data.Save(c.con.Request, c.con.Writer)
    return
}

func GetCookie(con *gin.Context) (cook *Cookie, err error) {
    data, err := Store.New(con.Request, "blckur")
    if err != nil {
        err = err.(securecookie.MultiError)[0]

        switch err {
        case securecookie.ErrMacInvalid:
            err = nil
        default:
            err = &errortypes.UnknownError{
                errors.Wrap(err, "auth: Unknown cookie error"),
            }
            return
        }
    }

    cook = &Cookie{
        data: data,
        con: con,
    }

    return
}

func Init() (err error) {
    db := database.GetDatabase()

    keyInt, err := settings.Get(db, "system", "cookie_key")
    if err != nil {
        err = &errortypes.UnknownError{
            errors.Wrap(err, "auth: Unknown error"),
        }
        return
    }

    var key []byte

    if keyInt == nil {
        key = securecookie.GenerateRandomKey(64)
        settings.Set(db, "system", "cookie_key", key)
    } else {
        key = keyInt.([]byte)
    }

    Store = sessions.NewCookieStore(key)

    return
}
