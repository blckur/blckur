package session

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	curCookieKey []byte
	Store        *sessions.CookieStore
)

type Cookie struct {
	Id    bson.ObjectId
	store *sessions.Session
	con   *gin.Context
}

func (c *Cookie) Get(key string) (val string) {
	valInter := c.store.Values[key]
	if valInter == nil {
		val = ""
	} else {
		val = valInter.(string)
	}
	return
}

func (c *Cookie) Set(key string, val string) {
	c.store.Values[key] = val
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
			err = &NotFoundError{
				errors.New("auth: Session not found"),
			}
		default:
			err = &errortypes.UnknownError{
				errors.Wrap(err, "auth: Unknown session error"),
			}
		}
		return
	}

	return
}

func (c *Cookie) NewSession(db *database.Database, id bson.ObjectId,
	remember bool) (sess *Session, err error) {

	sess, err = NewSession(db, id)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "auth: Unknown session error"),
		}
		return
	}

	c.Set("id", sess.Id.Hex())
	maxAge := 0

	if remember {
		maxAge = 15778500
	}

	c.store.Options = &sessions.Options{
		MaxAge: maxAge,
	}

	err = c.Save()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "auth: Unknown session error"),
		}
		return
	}

	return
}

func (c *Cookie) Remove(db *database.Database) (err error) {
	sessId := c.Get("id")
	if sessId == "" {
		return
	}

	err = RemoveSession(db, bson.ObjectIdHex(sessId))
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "auth: Unknown session error"),
		}
		return
	}

	c.Set("id", "")
	err = c.Save()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "auth: Unknown session error"),
		}
		return
	}

	return
}

func (c *Cookie) Save() (err error) {
	err = c.store.Save(c.con.Request, c.con.Writer)
	return
}

func GetCookie(con *gin.Context) (cook *Cookie, err error) {
	store, err := Store.New(con.Request, "blckur")
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
		store: store,
		con:   con,
	}

	return
}

func commitCookieKey() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	err = settings.Commit(db, settings.System, set.NewSet("cookie_key"))
	if err != nil {
		return
	}

	return
}

func update() {
	cookieKey := settings.System.CookieKey

	if len(cookieKey) == 0 {
		cookieKey = securecookie.GenerateRandomKey(64)
		settings.System.CookieKey = cookieKey

		for {
			err := commitCookieKey()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("session: Init cookie key")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}
	}

	if !bytes.Equal(curCookieKey, cookieKey) {
		Store = sessions.NewCookieStore(cookieKey)
		curCookieKey = cookieKey
	}
}

func init() {
	module := requires.New("session")
	module.After("settings")

	module.Handler = func() {
		update()
	}
}
