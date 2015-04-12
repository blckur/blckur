package cache

import (
	"container/list"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/dropbox/godropbox/errors"
	"github.com/garyburd/redigo/redis"
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

type message struct {
	Id   bson.ObjectId `json:"i"`
	Data string        `json:"d"`
}

type pubSubConn struct {
	address   string
	addMutex  *sync.Mutex
	remMutex  *sync.Mutex
	addQueue  *list.List
	remQueue  *list.List
	conn      *redis.PubSubConn
	listeners map[string]map[int]func(string)
	closed    bool
	counter   int
}

func (p *pubSubConn) Subscribe(key string, handler func(string)) (id int) {
	p.addMutex.Lock()

	handlers, ok := p.listeners[key]
	if !ok {
		handlers = map[int]func(string){}
		p.listeners[key] = handlers
		p.addQueue.PushBack(key)
	}

	id = p.counter
	p.counter += 1

	handlers[id] = handler

	p.addMutex.Unlock()

	return
}

func (p *pubSubConn) Unsubsribe(key string, id int) {
	p.remMutex.Lock()

	if handlers, ok := p.listeners[key]; ok {
		delete(handlers, id)
		if len(handlers) == 0 {
			delete(p.listeners, key)
			p.remQueue.PushBack(key)
		}
	}

	p.remMutex.Unlock()
}

func (p *pubSubConn) parseAddQueue() {
	for {
		if p.closed {
			return
		}

		p.addMutex.Lock()
		elem := p.addQueue.Front()
		conn := p.conn
		p.addMutex.Unlock()

		if elem == nil || conn == nil {
			break
		}

		err := conn.Subscribe(elem.Value.(string))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cache.pubsub: Subscribe error")
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		p.addMutex.Lock()
		p.addQueue.Remove(elem)
		p.addMutex.Unlock()
	}
}

func (p *pubSubConn) parseRemQueue() {
	for {
		if p.closed {
			return
		}

		p.remMutex.Lock()
		elem := p.remQueue.Front()
		conn := p.conn
		p.remMutex.Unlock()

		if elem == nil || conn == nil {
			break
		}

		err := conn.Unsubscribe(elem.Value.(string))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("cache.pubsub: Unsubscribe error")
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		p.remMutex.Lock()
		p.remQueue.Remove(elem)
		p.remMutex.Unlock()
	}
}

func (p *pubSubConn) Listen() {
	go func() {
		for {
			p.parseAddQueue()
			p.parseRemQueue()

			if p.closed {
				return
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	for {
		conn := p.conn
		if conn != nil {
		Loop:
			for {
				switch obj := conn.Receive().(type) {
				case redis.Message:
					if listeners, ok := p.listeners[obj.Channel]; ok {
						go func() {
							for _, listener := range listeners {
								listener(string(obj.Data))
							}
						}()
					}
				case error:
					if p.closed {
						return
					}

					logrus.WithFields(logrus.Fields{
						"error": obj,
					}).Error("cache.pubsub: Listen error")

					break Loop
				}
			}
		}

		for {
			if p.closed {
				return
			}

			time.Sleep(constants.RETRY_DELAY)

			conn, err := dialLong(p.address)
			if err != nil {
				err = &CacheError{
					errors.Wrap(err, "cache.pubsub: Dial error"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("cache.pubsub: Dial error")
				continue
			}

			p.addMutex.Lock()
			p.remMutex.Lock()

			p.conn = &redis.PubSubConn{
				Conn: conn,
			}

			p.addQueue = list.New()
			p.remQueue = list.New()

			for key, _ := range p.listeners {
				p.addQueue.PushBack(key)
			}

			p.addMutex.Unlock()
			p.remMutex.Unlock()

			break
		}
	}
}

func (p *pubSubConn) Close() {
	p.addMutex.Lock()
	p.remMutex.Lock()
	p.closed = true

	conn := p.conn
	if conn != nil {
		p.conn.Close()
		p.conn = nil
	}

	p.addMutex.Unlock()
	p.remMutex.Unlock()
}

func newPubSubConn(address string) (psc *pubSubConn) {
	psc = &pubSubConn{
		address:   address,
		addMutex:  &sync.Mutex{},
		remMutex:  &sync.Mutex{},
		addQueue:  list.New(),
		remQueue:  list.New(),
		listeners: map[string]map[int]func(string){},
		closed:    false,
	}

	conn, err := dialLong(address)
	if err != nil {
		err = &CacheError{
			errors.Wrap(err, "cache.pubsub: Connection error"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("cache.pubsub: Connection error")
	} else {
		psc.conn = &redis.PubSubConn{
			Conn: conn,
		}
	}

	return
}
