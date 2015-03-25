package cmd

import (
	"github.com/blckur/blckur/nodes"
)

// Starts clustered beanstalkd server and registers with database
func Queue() {
	opts := getServiceOptions()

	node := nodes.QueueNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
