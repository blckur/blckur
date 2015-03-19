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
	case "start":
		cmd.Server()
	case "beanstalkd":
		cmd.Beanstalkd()
	case "redis":
		cmd.Redis()
	}
}
