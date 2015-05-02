package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/handlers"
	"github.com/blckur/blckur/proxy"
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
	Debug      bool
}

func (a *AppNode) Start() {
	constants.Node = "app"

	port := strconv.Itoa(a.Port)

	router := gin.New()

	if a.Debug {
		router.Use(gin.Logger())
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

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
		"public_port": a.PublicPort,
		"version":     constants.Version,
		"debug":       a.Debug,
	}).Info("nodes.app: Starting app node")

	err := proxy.ListenAndServe(server)
	if err != nil {
		panic(err)
	}
}
