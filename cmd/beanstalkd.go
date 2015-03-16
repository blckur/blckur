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
	database.Init()
	settings.Init()
	logger.Init()

	id := os.Getenv("ID")
	if id == "" {
		panic("Node ID is required")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

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
