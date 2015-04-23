package queue

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/constants"
	"time"
)

type Stream struct {
	server string
	Stop   bool
}

func (q *Stream) Reserve(timeout time.Duration) (job *Job) {
	cacheConn := cache.Get()
	start := time.Time{}

	for {
		if q.Stop {
			job = nil
			return
		}

		conn, err := clst.conn(q.server)
		if err != nil {
			if start.IsZero() {
				start = time.Now()
			} else if time.Since(start) > constants.ErrLogDelay {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("queue.stream: Connection error")
			}

			time.Sleep(constants.RetryDelay)
			continue
		}

		id, body, err := conn.Reserve(timeout)
		if err != nil {
			if err.Error() == "reserve-with-timeout: timeout" ||
				err.Error() == "reserve-with-timeout: deadline soon" {
				continue
			}

			if start.IsZero() {
				start = time.Now()
			} else if time.Since(start) > constants.ErrLogDelay {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("queue.stream: Reserve error")
			}

			clst.close(q.server)
			continue
		}

		job = &Job{}
		err = json.Unmarshal(body, job)
		if err != nil {
			if start.IsZero() {
				start = time.Now()
			} else if time.Since(start) > constants.ErrLogDelay {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("queue.stream: Parse error")
			}

			continue
		}

		job.beanId = id
		job.conn = conn
		job.cacheConn = cacheConn

		val, err := cacheConn.GetString(job.Id.Hex())
		if err != nil {
			if start.IsZero() {
				start = time.Now()
			} else if time.Since(start) > constants.ErrLogDelay {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("queue.stream: Cache get error")
			}

			cacheConn = cache.Get()
			continue
		}

		if !start.IsZero() {
			start = time.Time{}
		}

		if val == "t" || time.Now().After(job.Timestamp.Add(job.Ttl)) {
			job.Delete()
			continue
		}

		break
	}

	return
}

func newStream(server string) (stream *Stream) {
	stream = &Stream{
		server: server,
		Stop:   false,
	}
	return
}
