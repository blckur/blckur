// Logger with output to stderr and papertrail.
package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/colorize"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"net"
	"os"
	"time"
)

var (
	buffer = make(chan *logrus.Entry, 128)
)

func formatLevel(lvl logrus.Level) (str string) {
	var colorBg colorize.Color

	switch lvl {
	case logrus.InfoLevel:
		colorBg = colorize.CyanBg
		str = "[INFO]"
	case logrus.WarnLevel:
		colorBg = colorize.YellowBg
		str = "[WARN]"
	case logrus.ErrorLevel:
		colorBg = colorize.RedBg
		str = "[ERRO]"
	case logrus.FatalLevel:
		colorBg = colorize.RedBg
		str = "[FATL]"
	case logrus.PanicLevel:
		colorBg = colorize.RedBg
		str = "[PANC]"
	default:
		colorBg = colorize.BlackBg
	}

	str = colorize.ColorString(str, colorize.WhiteBold, colorBg)

	return
}

func initSender() {
	var conn net.Conn

	go func() {
		go func() {
			conn = paperTrailConn()
		}()
		time.Sleep(2 * time.Second)

		for {
			entry := <-buffer

			if entry.Message[:6] == "logger" {
				continue
			}

			if conn != nil {
				err := paperTrailSend(conn, entry)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("logger: Papertrail error")
					conn = paperTrailConn()
				}
			}
		}
	}()

	messenger.Register("settings", "papertrail", func(_ *messenger.Message) {
		if conn != nil {
			conn.Close()
		}
	})
}

func init() {
	module := requires.New("logger")
	module.Before("database")

	module.Handler = func() {
		initSender()

		logrus.SetFormatter(&formatter{})
		logrus.AddHook(&logHook{})
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(logrus.InfoLevel)
	}
}
