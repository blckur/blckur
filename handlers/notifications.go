package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
)

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
