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
	port := strconv.Itoa(a.Port)
	address := getAddress() + ":" + port

	router := gin.Default()
	handlers.Register(router)

	server := &http.Server{
		Addr: a.Host + ":" + port,
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
