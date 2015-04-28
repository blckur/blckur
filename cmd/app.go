package cmd

import (
	"github.com/blckur/blckur/nodes"
	"os"
	"strconv"
)

// Starts web server node with JSON api
func App() {
	opts := getServiceOptions()

	source := os.Getenv("SOURCE")
	if source == "" {
		panic("Missing source")
	}

	var debug bool
	debugStr := os.Getenv("DEBUG")
	if debugStr == "" {
		debug = true
	} else {
		debug, _ = strconv.ParseBool(debugStr)
	}

	node := nodes.AppNode{
		Id:         opts.Id,
		Host:       opts.Host,
		Address:    opts.Address,
		Port:       opts.Port,
		PublicPort: opts.PublicPort,
		Source:     source,
		Debug:      debug,
	}

	node.Start()
}
