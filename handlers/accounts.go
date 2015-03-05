package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
	"math/rand"
)

type Account struct {
	Id bson.ObjectId `json:"id" binding:"required"`
	UserId bson.ObjectId `json:"user_id" binding:"required"`
	Type string `json:"type" binding:"required"`
	Identity string `json:"identity" binding:"required"`
}

type AccountPostData struct {
	Redirect string `json:"redirect" binding:"required"`
}

func accountsPost(c *gin.Context) {
	acct := &AccountPostData{}
	acctType := c.Params.ByName("type")
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	if acctType == "twitter" {
		url, err := account.AuthTwitter(db, sess.UserId)
		if err != nil {
			panic(err)
		}
		acct.Redirect = url
	} else {
		c.JSON(400, &ErrorData{
			Error: "unknown_type",
			Message: "Unknown account type",
		})
		return
	}

	c.JSON(200, acct)
}

func accountsGet(c *gin.Context) {
	accounts := []*Account{
		&Account{
			Id: bson.NewObjectId(),
			UserId: bson.NewObjectId(),
			Type: "gmail",
			Identity: "zach.huff.386@gmail.com",
		},
		&Account{
			Id: bson.NewObjectId(),
			UserId: bson.NewObjectId(),
			Type: "twitter",
			Identity: "@zachhuff386",
		},
		&Account{
			Id: bson.NewObjectId(),
			UserId: bson.NewObjectId(),
			Type: "github",
			Identity: "zachhuff386",
		},
	}

	c.JSON(200, accounts)
}
