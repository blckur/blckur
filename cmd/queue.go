package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts clustered beanstalkd server and registers with database
func Queue() {
	opts := getServiceOptions()

	node := nodes.QueueNode{
		Id:         opts.Id,
		Host:       opts.Host,
		Address:    opts.Address,
		Port:       opts.Port,
		PublicPort: opts.PublicPort,
	}

	node.Start()
}
