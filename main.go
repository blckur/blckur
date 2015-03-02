package main

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/handlers"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"flag"
)

func main() {
	flag.Parse()

	err := database.Init()
	if err != nil {
		panic(err)
	}

	err = settings.Init()
	if err != nil {
		panic(err)
	}

	switch flag.Arg(0) {
	case "set":
		group := flag.Arg(1)
		key := flag.Arg(2)
		val := flag.Arg(3)
		db := database.GetDatabase()
		defer db.Close()

		err = settings.Set(db, group, key, val)
		if err != nil {
			panic(err)
		}
		return
	case "start":
		err = session.Init()
		if err != nil {
			panic(err)
		}

		router := gin.Default()
		handlers.Register(router)

		server := &http.Server{
			Addr: ":80",
			Handler: router,
			ReadTimeout: 10 * time.Second,
			WriteTimeout: 10 * time.Second,
			MaxHeaderBytes: 4096,
		}

		server.ListenAndServe()
	}
}
