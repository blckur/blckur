package nodes

import (
	"github.com/blckur/blckur/handlers"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"strconv"
)

type AppNode struct {
	Id string
	Host string
	Port int
}

func (a *AppNode) Start() {
	router := gin.Default()
	handlers.Register(router)
	address := a.Host + ":" + strconv.Itoa(a.Port)

	server := &http.Server{
		Addr: address,
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"id": a.Id,
		"address": address,
	}).Info("nodes.app: Starting app node")

	server.ListenAndServe()
}
