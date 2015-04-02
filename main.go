package main

import (
	"github.com/blckur/blckur/cmd"
	"github.com/blckur/blckur/gdefer"
	"github.com/blckur/blckur/requires"
	"flag"
)

func main() {
	requires.Init()
	defer gdefer.End()
	flag.Parse()

	switch flag.Arg(0) {
	case "clear":
		cmd.Clear()
	case "set":
		cmd.Settings()
	case "app":
		cmd.App()
	case "queue":
		cmd.Queue()
	case "cache":
		cmd.Cache()
	case "scheduler":
		cmd.Scheduler()
	case "worker":
		cmd.Worker()
	case "pub":
		cmd.Pub()
	case "sub":
		cmd.Sub()
	case "pubsub":
		cmd.PubSub()
	case "xpub":
		cmd.StressPub()
	case "xsub":
		cmd.StressSub()
	}
}
