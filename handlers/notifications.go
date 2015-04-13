package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/session"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type notificationData struct {
	Id   bson.ObjectId `json:"id"`
	Read bool          `json:"read"`
}

func notificationGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	accts, err := notification.GetNotifications(db, sess.UserId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, accts)
}

func notificationPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	notfId := bson.ObjectIdHex(c.Params.ByName("id"))
	data := &notificationData{}

	c.Bind(data)

	notf, err := notification.GetNotification(db, sess.UserId, notfId)
	if err != nil {
		c.Fail(500, err)
		return
	}

	notf.Read = data.Read

	err = notf.CommitFields(db, set.NewSet("read"))
	if err != nil {
		c.Fail(500, err)
		return
	}

	pub := notification.NewPublisher(sess.UserId.Hex())
	defer pub.Close()

	// Ignore err
	_ = pub.Update(data)

	c.JSON(200, notf)
}
