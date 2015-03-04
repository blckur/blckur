package main

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/handlers"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/account"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"flag"
	"os"
)

func main() {
	utils.SeedRand()
	flag.Parse()

	database.Init()
	settings.Init()
	logger.Init()

	account.Init()
	messenger.Init()

	switch flag.Arg(0) {
	case "set":
		group := flag.Arg(1)
		key := flag.Arg(2)
		val := flag.Arg(3)
		db := database.GetDatabase()
		defer db.Close()

		err := settings.Set(db, group, key, val)
		if err != nil {
			panic(err)
		}

		err = messenger.Publish(db, "settings", group)
		if err != nil {
			panic(err)
		}

		return
	case "start":
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")
		if port == "" {
			port = "80"
		}
		addr := host + ":" + port
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
		}).Info("Starting server")

		server.ListenAndServe()
	}
}
