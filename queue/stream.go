package queue

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"time"
)

type stream struct {
	server string
	cluster *cluster
}

func (q *stream) Reserve(timeout time.Duration) (job *Job) {
	for {
		conn, err := q.cluster.conn(q.server)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("queue.stream: Connection error")
			continue
		}

		id, body, err := conn.Reserve(timeout)
		if err != nil {
			if err.Error() == "reserve-with-timeout: timeout" {
				continue
			}
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("queue.stream: Reserve error")
			q.cluster.close(q.server)
			continue
		}

		jobData := &JobData{}
		err = json.Unmarshal(body, jobData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("queue.stream: Parse error")
			continue
		}

		job = &Job{
			beanId: id,
			conn: conn,
			Id: jobData.Id,
			Type: jobData.Type,
			Data: jobData.Data,
		}

		break
	}

	return
}
