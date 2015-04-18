// Commands available in cli.
package cmd

import (
	"github.com/blckur/blckur/utils"
	"math/rand"
	"os"
	"strconv"
	"github.com/blckur/blckur/database"
)

type serviceOptions struct {
	Id   string
	Host string
	Port int
}

func getServiceOptions() (opts *serviceOptions) {
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

	database.MongoUrl = os.Getenv("DB")

	opts = &serviceOptions{
		Id:   id,
		Host: host,
		Port: port,
	}

	return
}
