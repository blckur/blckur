package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/dropbox/godropbox/errors"
	"github.com/blckur/blckur/utils"
)

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
				conn = clst.serverMap[server].Get()
				c.conns[server] = conn
			}

			v, e := redis.String(conn.Do("GET", key))
			if e != nil && e != redis.ErrNil {
				er = &CacheError{
					errors.Wrap(e, "cache: Get error"),
				}
				wait.Done()
				return
			}

			val = v
			success = true
			wait.Cancel()
		}(server)
	}

	wait.Wait()
	if !success {
		err = er
	}
	return
}
