package cmd

import (
	"github.com/blckur/blckur/utils"
	"os"
	"strconv"
	"math/rand"
)

type ServiceOptions struct {
	Id string
	Host string
	Port int
}

func GetServiceOptions() (opts *ServiceOptions) {
	id := os.Getenv("ID")
	if id == "" {
		id = utils.RandName()
	}

	host := os.Getenv("HOST")

	var port int
	portStr := os.Getenv("PORT")
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			panic(err)
		}
		port = p
	} else {
		port = rand.Intn(55000) + 10000
	}

	opts = &ServiceOptions{
		Id: id,
		Host: host,
		Port: port,
	}

	return
}
