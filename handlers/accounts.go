package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type AccountPostData struct {
	Redirect string `json:"redirect" binding:"required"`
}

type AccountData struct {
	Events []*account.EventType `json:"events"`
}

func accountsPost(c *gin.Context) {
	acct := &AccountPostData{}
	acctType := c.Params.ByName("type")
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	switch acctType {
	case "gmail":
		url, err := account.ReqGmail(db, sess.UserId)
		if err != nil {
			panic(err)
		}
		acct.Redirect = url
	case "twitter":
		url, err := account.ReqTwitter(db, sess.UserId)
		if err != nil {
			panic(err)
		}
		acct.Redirect = url
	default:
		c.JSON(400, &ErrorData{
			Error: "unknown_type",
			Message: "Unknown account type",
		})
		return
	}

	c.JSON(200, acct)
}

func accountsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	accts, err := account.GetAccounts(db, sess.UserId)
	if err != nil {
		panic(err)
	}

	c.JSON(200, accts)
}

func accountsPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	acctId := bson.ObjectIdHex(c.Params.ByName("id"))
	data := AccountData{}

	if !c.Bind(&data) {
		return
	}

	acct, err := account.GetAccount(db, sess.UserId, acctId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	acct.EventsParsed = data.Events
	account.UnmarshalGmail(acct)

	err = acct.CommitFields(set.NewSet("events"))
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, acct)
}
