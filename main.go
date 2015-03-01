package main

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/handlers"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	err := database.Init()
	if err != nil {
		panic(err)
	}

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
