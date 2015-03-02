package queue

import (
	"encoding/json"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"time"
	"log"
	"github.com/blckur/blckur/utils"
	"sync"
	"fmt"
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

	q.pool[server] = nil

	return
}

func (q *Queue) Conn() {
	for _, server := range q.servers {
		q.conn(server)
	}
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

func (q *Queue) marhsalJobs(data interface{}) (
		normalJob []byte, checkJob []byte, err error) {
	jobData := &JobData{
		Id: bson.NewObjectId(),
		Type: NORMAL,
		Data: data,
	}

	normalJob, err = json.Marshal(jobData)
	if err != nil {
		return
	}

	jobData.Type = CHECK

	checkJob, err = json.Marshal(jobData)
	if err != nil {
		return
	}

	return
}

func (q *Queue) putRetry(server string, data []byte, priority uint32,
		delay time.Duration, ttr time.Duration) (err error) {
	for i := 0; i < 2; i++ {
		conn, e := q.conn(server)
		if e != nil {
			err = e
			return
		}

		_, err = conn.Put(data, priority, delay, ttr)
		if err != nil {
			q.close(server)
		} else {
			break
		}
	}

	return
}

func (q *Queue) Put(data interface{}, priority uint32,
		delay time.Duration, ttr time.Duration) (err error) {
	normalJob, checkJob, err := q.marhsalJobs(data)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	servers := utils.NewStringStack(utils.ShuffleStringsNew(q.servers))

	waiters := &sync.WaitGroup{}
	waiters.Add(q.consistency)
	sent := 0
	sentMutex := &sync.Mutex{}

	for i := 0; i < q.consistency; i++ {
		go func(check bool) {
			for {
				server := servers.Pop()
				if server == "" {
					break
				}

				conn, err := q.conn(server)
				if err != nil {
					q.close(server)
					continue
				}

				if check {
					_, err = conn.Put(checkJob, priority,
						time.Duration(2) * ttr + delay, ttr)
					if err != nil {
						q.close(server)
						continue
					}
				} else {
					_, err = conn.Put(normalJob, priority, delay, ttr)
					if err != nil {
						q.close(server)
						continue
					}
				}

				sentMutex.Lock()
				sent += 1
				sentMutex.Unlock()

				waiters.Done()

				break
			}
		}(i == 0)
	}

	waiters.Wait()

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

		return
	}

	return
}

func (q *Queue) GetStreams() (streams []*QueueStream) {
	streams = make([]*QueueStream, len(q.servers))

	for i, server := range q.servers {
		streams[i] = &QueueStream{
			server: server,
			queue: q,
		}
	}

	return
}

func NewQueue(servers []string, consistency int) (que *Queue) {
	que = &Queue{
		servers: servers,
		consistency: consistency,
		pool: map[string]*beanstalk.Conn{},
	}
	return
}
