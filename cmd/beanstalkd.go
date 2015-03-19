package cmd

import (
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/logger"
	"os"
	"strconv"
)

func Beanstalkd() {
	logger.Init()
	database.Init()
	settings.Init()

	id := os.Getenv("ID")
	if id == "" {
		panic("Node ID is required")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "0"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	node := nodes.BeanstalkdNode{
		Id: id,
		Host: host,
		Port: portInt,
	}

	node.Start()
}
