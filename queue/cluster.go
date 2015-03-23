package queue

import (
	"encoding/json"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/stack"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"time"
	"sync"
	"fmt"
)

const (
	NORMAL = "normal"
	CHECK = "check"
)

type cluster struct {
	defaultConsistency int
	servers set.Set
	serversSlc []string
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
	for server := range c.servers.Iter() {
		c.conn(server.(string))
	}
}

func (c *cluster) Close() (err error) {
	for server := range c.servers.Iter() {
		err = c.close(server.(string))
		if err != nil {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "queue: Unknown error"),
			}
			return
		}
	}

	return
}

func (c *cluster) marhsalJob(job *Job) (
		normalJob []byte, checkJob []byte, err error) {
	job.queueType = NORMAL

	normalJob, err = json.Marshal(job)
	if err != nil {
		return
	}

	job.queueType = CHECK

	checkJob, err = json.Marshal(job)
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

func (c *cluster) Put(job *Job, priority int,
		delay time.Duration, ttr time.Duration) (err error) {
	normalJob, checkJob, err := c.marhsalJob(job)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown parse error"),
		}
		return
	}

	servers := stack.NewStringStack(utils.ShuffleStringsNew(c.serversSlc))

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
