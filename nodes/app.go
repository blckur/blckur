package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AppNode struct {
	Id         string
	Host       string
	Address    string
	Port       int
	PublicPort int
	Source     string
}

func (a *AppNode) Start() {
	port := strconv.Itoa(a.Port)

	router := gin.Default()
	handlers.Register(router, a.Source)

	server := &http.Server{
		Addr:           a.Host + ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"id":          a.Id,
		"address":     a.Address,
		"port":        a.Port,
		"public_port": a.Port,
	}).Info("nodes.app: Starting app node")

	server.ListenAndServe()
}
