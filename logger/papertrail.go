package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"net"
	"time"
)

func init() {
	senders = append(senders, &paperTrailSender{})
}

type paperTrailSender struct {
	limit  limiter
	buffer chan *logrus.Entry
	conn   net.Conn
}

func (p *paperTrailSender) listen() {
	p.conn = paperTrailConn()

	for {
		entry := <-p.buffer

		err := paperTrailSend(p.conn, entry)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("logger: Papertrail send error")
			p.conn = paperTrailConn()
		}
	}
}

func (p *paperTrailSender) Init() {
	p.limit = limiter{}
	p.buffer = make(chan *logrus.Entry, 128)

	go p.listen()

	messenger.Register("settings", "papertrail", func(_ *messenger.Message) {
		if p.conn != nil {
			p.conn.Close()
		}
	})
}

func (p *paperTrailSender) Parse(entry *logrus.Entry) {
	if settings.PapperTrail.Address == "" {
		return
	}

	if !p.limit.Check(entry,
		time.Duration(settings.PapperTrail.RateLimit)*time.Second) {

		return
	}

	if len(buffer) <= 125 {
		p.buffer <- entry
	}
}

func paperTrailConn() (conn net.Conn) {
	for {
		if settings.PapperTrail.Address != "" {
			c, err := net.Dial("udp", settings.PapperTrail.Address)
			if err != nil {
				err = &errortypes.UnknownError{
					errors.Wrap(err, "logger.papertrail: Connection error"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("logger.papertrail: Connection error")
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
			errors.Wrap(err, "logger.papertrail: Write error"),
		}
	}

	return
}
