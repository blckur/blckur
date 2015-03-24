package cmd

import (
	"github.com/blckur/blckur/nodes"
	"github.com/blckur/blckur/requires"
)

func Cache() {
	requires.Init()
	opts := GetServiceOptions()

	node := nodes.CacheNode{
		Id: opts.Id,
		Host: opts.Host,
		Port: opts.Port,
	}

	node.Start()
}
