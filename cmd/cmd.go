// Commands available in cli.
package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/utils"
	"math/rand"
	"os"
	"strconv"
)

type serviceOptions struct {
	Id         string
	Host       string
	Address    string
	Port       int
	PublicPort int
}

func getAddress() (address string) {
	address = os.Getenv("IP")
	if address != "" {
		return
	}

	address, err := utils.GetLocalAddress()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("cmd: Failed to get ip")
		panic(err)
	}

	return
}

func getPort() (port int) {
	portStr := os.Getenv("PORT")

	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cmd: Failed to parse port")
			panic(err)
		}
		port = p
	} else {
		port = rand.Intn(55000) + 10000
	}

	return
}

func getPublicPort() (port int) {
	dockerEndpoint := os.Getenv("DOCKER_API")

	if dockerEndpoint != "" {
		p, err := utils.GetDockerPort(dockerEndpoint)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cmd: Failed to get docker port")
			panic(err)
		}
		port = p
	}

	return
}

func getServiceOptions() (opts *serviceOptions) {
	id := os.Getenv("ID")
	if id == "" {
		id = utils.RandName()
	}

	host := os.Getenv("HOST")
	address := getAddress()
	port := getPort()

	pubPort := getPublicPort()
	if pubPort == 0 {
		pubPort = port
	}

	opts = &serviceOptions{
		Id:         id,
		Host:       host,
		Address:    address,
		Port:       port,
		PublicPort: pubPort,
	}

	return
}
