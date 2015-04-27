package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"net"
	"time"
)

func paperTrailConn() (conn net.Conn) {
	for {
		if settings.PapperTrail.Address != "" {
			c, err := net.Dial("udp", settings.PapperTrail.Address)
			if err != nil {
				err = &errortypes.UnknownError{
					errors.Wrap(err, "logger: Papertrail connection"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("logger: Papertrail connection")
			} else {
				conn = c
				break
			}
		}

		time.Sleep(constants.RetryDelay)
	}

	return
}

func paperTrailSend(conn net.Conn, entry *logrus.Entry) (err error) {
	msg, err := entry.String()
	if err != nil {
		return
	}

	_, err = conn.Write([]byte(msg))
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "logger: Papertrail write"),
		}
	}

	return
}
