package auth

import (
    "github.com/blckur/blckur/database"
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
    val := c.Get("id")
    if val == "" {
        err = &NotFoundError{
            errors.New("auth: Session not found"),
        }
        return
    }

    sess, err = GetSession(db, bson.ObjectIdHex(val))
    if err != nil {
        err = &UnknownError{
            errors.Wrap(err, "auth: Unknown session error"),
        }
        return
    }

    return
}

func (c *Cookie) NewSession(db *database.Database, id bson.ObjectId) (
        sess *Session, err error) {
    sess, err = NewSession(db, id)
    if err != nil {
        err = &UnknownError{
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
            err = &UnknownError{
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

func Init() {
    Store = sessions.NewCookieStore([]byte("todo"))
}
