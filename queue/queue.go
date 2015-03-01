package queue

import (
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"time"
//	"log"
)

var (
	NORMAL = "normal"
	CHECK = "check"
)

type Job struct {
	beanId uint64
	conn *beanstalk.Conn
	Id bson.ObjectId
	Type string
	Data interface{}
}

func (j *Job) Delete() (err error) {
	err = j.conn.Delete(j.beanId)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown error"),
		}
		return
	}

	return
}

type JobData struct {
	Id bson.ObjectId `json:"id"`
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

type Queue struct {
	consistency int
	servers []string
	pool map[string]*beanstalk.Conn
}

func (q *Queue) conn(server string) (conn *beanstalk.Conn, err error) {
	conn = q.pool[server]
	if conn != nil {
		return
	}

	conn, err = beanstalk.Dial("tcp", server)
	if err != nil {
		return
	}

	q.pool[server] = conn

	return
}

func (q *Queue) close(server string) (err error) {
	err = q.pool[server].Close()
	if err != nil {
		return
	}

	println("remove")
	q.pool[server] = nil

	return
}

func (q *Queue) Close() (err error) {
	for _, server := range q.servers {
		err = q.close(server)
		if err != nil {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "queue: Unknown error"),
			}
			return
		}
	}

	return
}

func (q *Queue) Put(data interface{}, priority uint32,
		delay time.Duration, ttr time.Duration) (err error) {
	jobId := bson.NewObjectId()

	normalJobData := &JobData{
		Id: jobId,
		Type: NORMAL,
		Data: data,
	}
	normalJob, err := json.Marshal(normalJobData)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	checkJobData := &JobData{
		Id: jobId,
		Type: CHECK,
		Data: data,
	}
	checkJob, err := json.Marshal(checkJobData)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	servers := make([]string, len(q.servers))
	copy(servers, q.servers)

	sent := 0
	for _, server := range servers {
		conn, err := q.conn(server)
		if err != nil {
			continue
		}

		if sent == 0 {
			_, err = conn.Put(normalJob, priority, delay, ttr)
			if err != nil {
				continue
			}
		} else {
			_, err = conn.Put(checkJob, priority,
				time.Duration(2) * ttr + delay, ttr)
			if err != nil {
				continue
			}
		}

		sent += 1

		if sent >= q.consistency {
			break
		}
	}

	if sent < q.consistency {
		if err != nil {
			err = &JobFailed{
				errors.Wrap(err, "queue: Job queue failed"),
			}
		} else {
			err = &JobFailed{
				errors.New("queue: Job queue failed"),
			}
		}
		return
	} else {
		err = nil
	}

	return
}

func (q *Queue) Reserve(timeout time.Duration) (job *Job) {
	for {
		for _, server := range q.servers {
			conn, err := q.conn(server)
			if err != nil {
				//log.Printf("ERROR: queue: %s", err.Error())
				continue
			}

			id, body, err := conn.Reserve(timeout)
			if err != nil {
				if err.Error() == "reserve-with-timeout: timeout" {
					continue
				}
				//log.Printf("ERROR: queue: %s", err.Error())
				q.close(server)
				continue
			}

			jobData := &JobData{}
			err = json.Unmarshal(body, jobData)
			if err != nil {
				//log.Printf("ERROR: queue: %s", err.Error())
				continue
			}

			job = &Job{
				beanId: id,
				conn: conn,
				Id: jobData.Id,
				Type: jobData.Type,
				Data: jobData.Data,
			}
			return
		}
	}
}

func NewQueue(servers []string, consistency int) (que *Queue) {
	que = &Queue{
		servers: servers,
		consistency: consistency,
		pool: map[string]*beanstalk.Conn{},
	}
	return
}
