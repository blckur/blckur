package main

import (
	"github.com/blckur/blckur/cmd"
	"flag"
)

func main() {
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
