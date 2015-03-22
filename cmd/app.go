package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/handlers"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/account"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"os"
)

func App() {
	logger.Init()
	database.Init()
	settings.Init()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	addr := host + ":" + port

	account.Init()
	cache.Init()
	messenger.Init()
	session.Init()

	router := gin.Default()
	handlers.Register(router)

	server := &http.Server{
		Addr: addr,
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"address": addr,
	}).Info("nodes.app: Starting app node")

	server.ListenAndServe()
}
