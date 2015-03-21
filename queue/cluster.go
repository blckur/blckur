package queue

import (
	"encoding/json"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/stack"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"time"
	"sync"
	"fmt"
)

var (
	NORMAL = "normal"
	CHECK = "check"
)

type JobData struct {
	Id bson.ObjectId `json:"id"`
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

type cluster struct {
	consistency int
	servers []string
	pool map[string]*beanstalk.Conn
}

func (q *cluster) conn(server string) (conn *beanstalk.Conn, err error) {
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

func (q *cluster) close(server string) (err error) {
	err = q.pool[server].Close()
	if err != nil {
		return
	}

	q.pool[server] = nil

	return
}

func (q *cluster) Conn() {
	for _, server := range q.servers {
		q.conn(server)
	}
}

func (q *cluster) Close() (err error) {
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

func (q *cluster) marhsalJobs(data interface{}) (
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

func (q *cluster) putRetry(server string, data []byte, priority uint32,
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

func (q *cluster) Put(data interface{}, priority uint32,
		delay time.Duration, ttr time.Duration) (err error) {
	normalJob, checkJob, err := q.marhsalJobs(data)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	servers := stack.NewStringStack(utils.ShuffleStringsNew(q.servers))

	waiters := &sync.WaitGroup{}
	waiters.Add(q.consistency)
	sent := 0
	sentMutex := sync.Mutex{}

	for i := 0; i < q.consistency; i++ {
		go func(normal bool) {
			var err error

			for {
				server := servers.Pop()
				if server == "" {
					break
				}

				if normal {
					err = q.putRetry(server, normalJob, priority, delay, ttr)
				} else {
					err = q.putRetry(server, checkJob, priority,
					time.Duration(2) * ttr + delay, ttr)
				}
				if err != nil {
					continue
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
		msg := fmt.Sprintf("queue: Job consistency unmet %d/%d",
		sent, q.consistency)

		if err != nil {
			err = &JobFailed{
				errors.Wrap(err, msg),
			}
		} else {
			err = &JobFailed{
				errors.New(msg),
			}
		}
		return
	} else {
		err = nil
	}

	return
}

func (q *cluster) GetStreams() (streams []*stream) {
	streams = make([]*stream, len(q.servers))

	for i, server := range q.servers {
		streams[i] = &stream{
			server: server,
			queue: q,
		}
	}

	return
}
