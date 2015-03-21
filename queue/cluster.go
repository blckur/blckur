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
	defaultConsistency int
	servers []string
	pool map[string]*beanstalk.Conn
}

func (c *cluster) conn(server string) (conn *beanstalk.Conn, err error) {
	conn = c.pool[server]
	if conn != nil {
		return
	}

	conn, err = beanstalk.Dial("tcp", server)
	if err != nil {
		return
	}

	c.pool[server] = conn

	return
}

func (c *cluster) close(server string) (err error) {
	err = c.pool[server].Close()
	if err != nil {
		return
	}

	c.pool[server] = nil

	return
}

func (c *cluster) Conn() {
	for _, server := range c.servers {
		c.conn(server)
	}
}

func (c *cluster) Close() (err error) {
	for _, server := range c.servers {
		err = c.close(server)
		if err != nil {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "queue: Unknown error"),
			}
			return
		}
	}

	return
}

func (c *cluster) marhsalJobs(data interface{}) (
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

func (c *cluster) putRetry(server string, data []byte, priority int,
		delay time.Duration, ttr time.Duration) (err error) {
	for i := 0; i < 2; i++ {
		conn, e := c.conn(server)
		if e != nil {
			err = e
			return
		}

		_, err = conn.Put(data, uint32(priority), delay, ttr)
		if err != nil {
			c.close(server)
		} else {
			break
		}
	}

	return
}

func (c *cluster) Put(data interface{}, priority int,
		delay time.Duration, ttr time.Duration) (err error) {
	normalJob, checkJob, err := c.marhsalJobs(data)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	servers := stack.NewStringStack(utils.ShuffleStringsNew(c.servers))

	waiters := &sync.WaitGroup{}
	waiters.Add(c.defaultConsistency)
	sent := 0
	sentMutex := sync.Mutex{}

	for i := 0; i < c.defaultConsistency; i++ {
		go func(normal bool) {
			var err error

			for {
				server := servers.Pop()
				if server == "" {
					break
				}

				if normal {
					err = c.putRetry(server, normalJob, priority, delay, ttr)
				} else {
					err = c.putRetry(server, checkJob, priority,
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

	if sent < c.defaultConsistency {
		msg := fmt.Sprintf("queue: Job consistency unmet %d/%d",
			sent, c.defaultConsistency)

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

func (c *cluster) GetStreams() (streams []*stream) {
	streams = make([]*stream, len(c.servers))

	for i, server := range c.servers {
		streams[i] = &stream{
			server: server,
			cluster: c,
		}
	}

	return
}
