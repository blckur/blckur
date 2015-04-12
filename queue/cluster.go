package queue

import (
	"encoding/json"
	"fmt"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/stack"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"sync"
	"time"
)

type cluster struct {
	defaultConsistency int
	servers            set.Set
	serversSlc         []string
	pool               map[string]*beanstalk.Conn
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
	srv, ok := c.pool[server]
	if !ok {
		return
	}

	err = srv.Close()
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
				errors.Wrap(err, "queue.cluster: Unknown error"),
			}
			return
		}
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

	jsonJob, err := json.Marshal(job)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue.cluster: Unknown parse error"),
		}
		return
	}

	servers := stack.NewStringStack(utils.ShuffleStringsNew(c.serversSlc))

	waiters := &sync.WaitGroup{}
	waiters.Add(c.defaultConsistency)
	sent := 0
	sentMutex := sync.Mutex{}

	for i := 0; i < c.defaultConsistency; i++ {
		go func(first bool) {
			var err error

			for {
				server := servers.Pop()
				if server == "" {
					break
				}

				if first {
					err = c.putRetry(server, jsonJob, priority, delay, ttr)
				} else {
					err = c.putRetry(server, jsonJob, priority,
						time.Duration(2)*ttr+delay, ttr)
				}
				if err != nil {
					continue
				}

				sentMutex.Lock()
				sent += 1
				sentMutex.Unlock()

				break
			}

			waiters.Done()
		}(i == 0)
	}

	waiters.Wait()

	if sent < c.defaultConsistency {
		msg := fmt.Sprintf("queue.cluster: Job consistency unmet %d/%d",
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
