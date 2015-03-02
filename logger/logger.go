package logger

import (
	"github.com/blckur/blckur/utils"
	"github.com/op/go-logging"
	"os"
)

var (
	log = logging.MustGetLogger("blckur")
	format = logging.MustStringFormatter(
		"%{color:bold}%{time:06/01/02 15:04:05 }" +
		"%{shortfunc} %{level:.4s} ▶ %{color:reset} %{message}")
)

func Init() (err error) {
	utils.After("database")

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatted := logging.NewBackendFormatter(backend, format)
	leveled := logging.AddModuleLevel(formatted)
	leveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(leveled)

	utils.Register("logger")

	return
}
