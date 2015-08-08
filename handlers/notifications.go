package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/user"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
	"time"
)

type notificationData struct {
	Id      bson.ObjectId `json:"id"`
	Read    bool          `json:"read"`
	Origin  string        `json:"origin"`
	Link    string        `json:"link"`
	Subject string        `json:"subject"`
	Body    string        `json:"body"`
}

func notificationGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	// TODO
	accts, err := notification.GetNotifications(db, sess.UserId, 32)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, accts)
}

func notificationDel(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	notfId := bson.ObjectIdHex(c.Params.ByName("id"))

	err := notification.RemNotification(db, sess.UserId, notfId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	pub := notification.NewPublisher(sess.UserId.Hex())
	pub.Remove(notfId.Hex())
	defer pub.Close()

	c.JSON(200, map[string]string{})
}

func notificationPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	notfId := bson.ObjectIdHex(c.Params.ByName("id"))
	data := &notificationData{}

	c.Bind(data)

	notf, err := notification.GetNotification(db, sess.UserId, notfId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	notf.Read = data.Read

	err = notf.CommitFields(db, set.NewSet("read"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	pub := notification.NewPublisher(sess.UserId.Hex())
	pub.Update(data)
	defer pub.Close()

	c.JSON(200, notf)
}

func notificationPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	apikey := c.Params.ByName("apikey")
	if apikey == "" {
		c.AbortWithStatus(401)
		return
	}

	usr, err := user.FindUserApiKey(db, apikey)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	data := &notificationData{}
	c.Bind(data)

	notf := notification.Notification{
		UserId:    usr.Id,
		AccountId: usr.Id,
		RemoteId:  bson.NewObjectId().String(),
		Timestamp: time.Now(),
		Type:      "blckur",
		Origin:    data.Origin,
		Link:      data.Link,
		Subject:   data.Subject,
		Body:      data.Body,
	}

	_, err = notf.Initialize(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	pub := notification.NewPublisher(usr.Id.Hex())
	defer pub.Close()
	pub.New(notf)

	c.JSON(200, notf)
}
