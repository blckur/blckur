package utils

import (
	"github.com/fsouza/go-dockerclient"
	"os"
	"strconv"
)

func GetDockerPort(endpoint string) (port int, err error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return
	}

	cont, err := client.InspectContainer(os.Getenv("HOSTNAME"))
	if err != nil {
		return
	}

	for _, bindings := range cont.NetworkSettings.Ports {
		portStr := bindings[0].HostPort
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return
		}
	}

	return
}
