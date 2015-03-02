package logger

import (
	"github.com/blckur/blckur/utils"
	"github.com/op/go-logging"
	"os"
)

var (
	log = logging.MustGetLogger("blckur")
	format = logging.MustStringFormatter(
		"%{color:bold}%{time:06/01/02 15:04:05} " +
		"%{level:.4s} ▶%{color:reset} %{message}")
)

func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}

func Notice(format string, args ...interface{}) {
	log.Notice(format, args...)
}

func Warning(format string, args ...interface{}) {
	log.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}

func Critical(format string, args ...interface{}) {
	log.Critical(format, args...)
}

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
