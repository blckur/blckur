package queue

import (
	"encoding/json"
	"time"
	"log"
)

type QueueStream struct {
	server string
	queue *Queue
}

func (q *QueueStream) Reserve(timeout time.Duration) (job *Job) {
	for {
		conn, err := q.queue.conn(q.server)
		if err != nil {
			log.Printf("ERROR: queue: %s", err.Error())
			continue
		}

		id, body, err := conn.Reserve(timeout)
		if err != nil {
			if err.Error() == "reserve-with-timeout: timeout" {
				continue
			}
			log.Printf("ERROR: queue: %s", err.Error())
			q.queue.close(q.server)
			continue
		}

		jobData := &JobData{}
		err = json.Unmarshal(body, jobData)
		if err != nil {
			log.Printf("ERROR: queue: %s", err.Error())
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
