package main

import (
	"github.com/blckur/blckur/cmd"
	"github.com/blckur/blckur/gdefer"
	"flag"
)

func main() {
	defer gdefer.End()
	flag.Parse()

	switch flag.Arg(0) {
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
	}
}
