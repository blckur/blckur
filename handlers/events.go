package handlers

import (
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

const (
	writeTimeout = 10 * time.Second
	pingInterval = 30 * time.Second
	pingWait     = 40 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 30 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO Check r.Header.Get("Origin")
			return true
		},
	}
)

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

	conn.SetReadDeadline(time.Now().Add(pingWait))
	conn.SetPongHandler(func(x string) (err error) {
		conn.SetReadDeadline(time.Now().Add(pingWait))
		return
	})

	lst := cache.Subscribe(sess.UserId.Hex())
	ticker := time.NewTicker(pingInterval)
	sub := lst.Listen()

	defer func() {
		ticker.Stop()
		lst.Close()
		conn.Close()
	}()

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}
		}
	}()

	for {
		select {
		case msg, ok := <-sub:
			if !ok {
				conn.WriteControl(websocket.CloseMessage, []byte{},
					time.Now().Add(writeTimeout))
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err = conn.WriteJSON(msg)
			if err != nil {
				return
			}

		case <-ticker.C:
			err = conn.WriteControl(websocket.PingMessage, []byte{},
				time.Now().Add(writeTimeout))
			if err != nil {
				return
			}
		}
	}
}
