package handlers

import (
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"labix.org/v2/mgo/bson"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO Check r.Header.Get("Origin")
		return true
	},
}

type event struct {
	Id       bson.ObjectId `json:"id"`
	Type     string        `json:"type"`
	Resource string        `json:"resource_id"`
}

func eventGet(c *gin.Context) {
	sess := c.MustGet("session").(*session.Session)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Fail(500, err)
		return
	}

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}
		}
	}()

	lst := cache.Subscribe(sess.UserId.Hex())
	defer lst.Close()

	for msg := range lst.Listen() {
		evt := &event{
			Id:   bson.NewObjectId(),
			Type: msg,
		}

		err = conn.WriteJSON(evt)
		if err != nil {
			c.Fail(500, err)
			return
		}
	}

	c.Fail(500, nil)
}
