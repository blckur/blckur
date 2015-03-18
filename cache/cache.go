package cache

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/shard"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/nodes"
	"github.com/garyburd/redigo/redis"
	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	clst *cluster
)

type cluster struct {
	serverMap map[string]*redis.Pool
	shrd *shard.Shard
}

func dial(node *nodes.Node) (conn redis.Conn, err error) {
	timeout := time.Duration(settings.Redis.TimeoutMilli) * time.Millisecond
	conn, err = redis.DialTimeout("tcp", node.Address, timeout, timeout,
		timeout)
	return
}

func newPool(node *nodes.Node) (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle: settings.Redis.MaxIdle,
		MaxActive: settings.Redis.MaxActive,
		IdleTimeout: time.Duration(settings.Redis.IdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = dial(node)
			return
		},
		TestOnBorrow: func(conn redis.Conn, _ time.Time) (err error) {
			_, err = conn.Do("PING")
			return
		},
	}
	return
}

func Get() (conn *ClusterConn) {
	conn = &ClusterConn{
		clst: clst,
		conns: map[string]redis.Conn{},
	}
	return
}

func update() {
	for {
		db := database.GetDatabase()
		coll := db.Nodes()
		nodes := []*nodes.Node{}
		cls := &cluster{
			serverMap: map[string]*redis.Pool{},
		}

		err := coll.Find(bson.M{
			"type": "redis",
		}).All(&nodes)
		if err != nil {
			err = database.ParseError(err)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cache: Database error")
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		svrKeys := []string{}

		for _, node := range nodes {
			svrKeys = append(svrKeys, node.Id)
			cls.serverMap[node.Id] = newPool(node)
		}

		cls.shrd = shard.New(svrKeys, settings.Redis.Consistency)

		clst = cls

		break
	}
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "redis", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("redis", "redis", func(_ *messenger.Message) {
		go update()
	})
	update()

	requires.Register("cache")
}
