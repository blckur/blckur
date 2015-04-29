// Redis cluster with client side sharding/replication.
//
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
//
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
//	conn.Close()
package cache

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/gdefer"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/node"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/shard"
	"github.com/dropbox/godropbox/container/set"
	"github.com/garyburd/redigo/redis"
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

var (
	clst      *cluster
	subs      set.Set
	clstMutex = sync.RWMutex{}
)

type cluster struct {
	serverMap   map[string]*redis.Pool
	pubsubConns map[string]*pubSubConn
	shrd        *shard.Shard
}

func dial(address string) (conn redis.Conn, err error) {
	timeout := time.Duration(settings.Cache.TimeoutMilli) * time.Millisecond
	conn, err = redis.DialTimeout("tcp", address, timeout, timeout, timeout)
	return
}

func dialLong(address string) (conn redis.Conn, err error) {
	timeout := time.Duration(settings.Cache.TimeoutMilli) * time.Millisecond
	timeout2 := time.Duration(0)
	conn, err = redis.DialTimeout("tcp", address, timeout, timeout2, timeout2)
	return
}

func newPool(address string) (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle:     settings.Cache.MaxIdle,
		MaxActive:   settings.Cache.MaxActive,
		IdleTimeout: time.Duration(settings.Cache.IdleTimeout) * time.Second,
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
	}
	return
}

func update() {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Nodes()
	nodes := []*node.Node{}
	cls := &cluster{
		serverMap:   map[string]*redis.Pool{},
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
		time.Sleep(constants.RetryDelay)
		update()
	}

	svrKeys := []string{}

	for _, node := range nodes {
		svrKeys = append(svrKeys, node.Id)
		cls.serverMap[node.Id] = newPool(node.Address)
		cls.pubsubConns[node.Id] = newPubSubConn(node.Address)
	}

	cls.shrd = shard.New(svrKeys, settings.Cache.Consistency)

	clstMutex.Lock()

	curClst := clst
	clst = cls

	if curClst != nil {
		for _, psc := range curClst.pubsubConns {
			psc.Close()
		}
	}

	for _, psc := range clst.pubsubConns {
		go func(psc *pubSubConn) {
			psc.Listen()
		}(psc)
	}

	for lstnrIntf := range subs.Iter() {
		lstnr := lstnrIntf.(*Listener)
		lstnr.reshard()
	}

	clstMutex.Unlock()
}

func init() {
	module := requires.New("cache")
	module.After("settings")
	module.Before("messenger")

	module.Handler = func() {
		messenger.Register("settings", "cache", func(_ *messenger.Message) {
			go update()
		})
		messenger.Register("cache", "update", func(_ *messenger.Message) {
			go update()
		})
		update()

		gdefer.Defer(func() {
			for _, pool := range clst.serverMap {
				pool.Close()
			}

			for _, psc := range clst.pubsubConns {
				psc.Close()
			}
		})
	}

	subs = set.NewSet()
}
