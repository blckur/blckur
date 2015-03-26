package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type accountPostData struct {
	Redirect string `json:"redirect"`
}

type accountData struct {
	Alerts []*account.Alert `json:"alerts"`
}

func accountsPost(c *gin.Context) {
	acct := &accountPostData{}
	acctType := c.Params.ByName("type")
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	auth, _, err := account.GetAuth(acctType)
	if err != nil {
		c.JSON(400, &errorData{
			Error: "unknown_type",
			Message: "Unknown account type",
		})
	}

	url, err := auth.Request(db, sess.UserId)
	if err != nil {
		c.Fail(500, err)
		return
	}
	acct.Redirect = url

	c.JSON(200, acct)
}

func accountsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	accts, err := account.GetAccounts(db, sess.UserId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, accts)
}

func accountsPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	acctId := bson.ObjectIdHex(c.Params.ByName("id"))
	data := accountData{}

	if !c.Bind(&data) {
		return
	}

	acct, err := account.GetAccount(db, sess.UserId, acctId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	acct.Alerts = data.Alerts
	acct.ParseAlerts()

	err = acct.CommitFields(db, set.NewSet("alerts"))
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, acct)
}

func accountsDel(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	acctId := bson.ObjectIdHex(c.Params.ByName("id"))

	err := account.RemAccount(db, sess.UserId, acctId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, map[string]string{})
}
