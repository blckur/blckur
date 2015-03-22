package cache

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/shard"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/gdefer"
	"github.com/blckur/blckur/nodes"
	"github.com/garyburd/redigo/redis"
	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

// Get set example
// 	conn := cache.Get()
// 	err := conn.SetString("key", "val")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	val, err := conn.GetString("key")
// 	if err != nil {
// 		panic(err)
// 	}
//	println(val)
//
//	conn.Close()

// Pub sub example
//	lst := cache.Subscribe("channel")
//
//	go func() {
//		for msg := range lst.Listen() {
//			println(msg)
//			lst.Close()
//		}
//	}()
//
//	time.Sleep(50 * time.Millisecond)
//
//	conn := cache.Get()
//	err := conn.Publish("channel", "message")
//	if err != nil {
//		panic(err)
//	}
//
// conn.Close()

var (
	clst *cluster
	clstMutex = sync.RWMutex{}
)

type cluster struct {
	serverMap map[string]*redis.Pool
	pubsubConns map[string]*pubSubConn
	shrd *shard.Shard
}

func dial(address string) (conn redis.Conn, err error) {
	timeout := time.Duration(settings.Redis.TimeoutMilli) * time.Millisecond
	conn, err = redis.DialTimeout("tcp", address, timeout, timeout, timeout)
	return
}

func dialLong(address string) (conn redis.Conn, err error) {
	timeout := time.Duration(settings.Redis.TimeoutMilli) * time.Millisecond
	timeout2 := time.Duration(0)
	conn, err = redis.DialTimeout("tcp", address, timeout, timeout2, timeout2)
	return
}

func newPool(address string) (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle: settings.Redis.MaxIdle,
		MaxActive: settings.Redis.MaxActive,
		IdleTimeout: time.Duration(settings.Redis.IdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = dial(address)
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
			pubsubConns: map[string]*pubSubConn{},
		}

		err := coll.Find(bson.M{
			"type": "cache",
		}).All(&nodes)
		if err != nil {
			err = database.ParseError(err)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cache: Update error")
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		svrKeys := []string{}

		for _, node := range nodes {
			svrKeys = append(svrKeys, node.Id)
			cls.serverMap[node.Id] = newPool(node.Address)
			cls.pubsubConns[node.Id] = newPubSubConn(node.Address)
		}

		cls.shrd = shard.New(svrKeys, settings.Redis.Consistency)

		clstMutex.Lock()

		curClst := clst
		clst = cls

		if curClst != nil {
			for _, psc := range curClst.pubsubConns {
				psc.Reshard()
			}
		}

		for _, psc := range clst.pubsubConns {
			go func(psc *pubSubConn) {
				psc.Listen()
			}(psc)
		}

		clstMutex.Unlock()

		break
	}
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "redis", func(_ *messenger.Message) {
		go update()
	})
	messenger.Register("cache", "update", func(_ *messenger.Message) {
		go update()
	})
	update()

	requires.Register("cache")

	gdefer.Defer(func() {
		for _, pool := range clst.serverMap {
			pool.Close()
		}

		for _, psc := range clst.pubsubConns {
			psc.Close()
		}
	})
}
