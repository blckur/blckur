package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/dropbox/godropbox/errors"
	"github.com/blckur/blckur/utils"
	"time"
	"strconv"
)

// Connection to shared cluster
type ClusterConn struct {
	clst *cluster
	conns map[string]redis.Conn
}

func (c *ClusterConn) GetString(key string) (val string, err error) {
	wait := utils.WaitCancel{}
	success := false
	var er error

	for _, server := range c.clst.shrd.Select(key) {
		wait.Add(1)
		go func(server string) {
			conn, ok := c.conns[server]
			if !ok {
				conn = c.clst.serverMap[server].Get()
				c.conns[server] = conn
			}

			v, e := redis.String(conn.Do("GET", key))
			if e != nil && e != redis.ErrNil {
				er = &CacheError{
					errors.Wrap(e, "cache.cluster: Get error"),
				}
				wait.Done()
				return
			}

			val = v
			success = true
			if e == redis.ErrNil {
				wait.Done()
			} else {
				wait.Cancel()
			}
		}(server)
	}

	wait.Wait()
	if !success {
		err = er
	}
	return
}

func (c *ClusterConn) SetString(key string, val string,
		ttl time.Duration) (err error) {
	wait := utils.WaitCancel{}
	success := false
	var er error

	for _, server := range c.clst.shrd.Select(key) {
		wait.Add(1)
		go func(server string) {
			conn, ok := c.conns[server]
			if !ok {
				conn = c.clst.serverMap[server].Get()
				c.conns[server] = conn
			}

			_, e := conn.Do("SETEX", key,
				strconv.FormatInt(int64(ttl.Seconds()), 10), val)
			if e != nil {
				er = &CacheError{
					errors.Wrap(e, "cache.cluster: Set error"),
				}
				wait.Done()
				return
			}

			success = true
			wait.Done()
		}(server)
	}

	wait.Wait()
	if !success {
		err = er
	}
	return
}

func (c *ClusterConn) Publish(channel string, val string) (err error) {
	wait := utils.WaitCancel{}
	success := false
	var er error

	for _, server := range c.clst.shrd.Select(channel) {
		wait.Add(1)
		go func(server string) {
			conn, ok := c.conns[server]
			if !ok {
				conn = c.clst.serverMap[server].Get()
				c.conns[server] = conn
			}

			_, e := conn.Do("PUBLISH", channel, val)
			if e != nil {
				er = &CacheError{
					errors.Wrap(e, "cache.cluster: Publish error"),
				}
				wait.Done()
				return
			}

			success = true
			wait.Done()
		}(server)
	}

	wait.Wait()
	if !success {
		err = er
	}
	return
}

func (c *ClusterConn) Close() {
	cns := c.conns
	c.conns = map[string]redis.Conn{}
	for _, conn := range cns {
		conn.Close()
	}
}
