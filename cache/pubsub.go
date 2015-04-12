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

type queueItem struct {
	State bool
	Key   string
}

type pubSubConn struct {
	address   string
	mutex     *sync.Mutex
	queue     *list.List
	conn      *redis.PubSubConn
	listeners map[string]map[int]func(string)
	closed    bool
	counter   int
}

func (p *pubSubConn) Subscribe(key string, handler func(string)) (id int) {
	p.mutex.Lock()

	handlers, ok := p.listeners[key]
	if !ok {
		handlers = map[int]func(string){}
		p.listeners[key] = handlers

		p.queue.PushBack(queueItem{
			State: true,
			Key:   key,
		})
	}

	id = p.counter
	p.counter += 1

	handlers[id] = handler

	p.mutex.Unlock()

	return
}

func (p *pubSubConn) Unsubsribe(key string, id int) {
	p.mutex.Lock()

	if handlers, ok := p.listeners[key]; ok {
		delete(handlers, id)
		if len(handlers) == 0 {
			delete(p.listeners, key)

			p.queue.PushBack(queueItem{
				State: false,
				Key:   key,
			})
		}
	}

	p.mutex.Unlock()
}

func (p *pubSubConn) parseQueue() {
	for {
		if p.closed {
			return
		}

		p.mutex.Lock()
		elem := p.queue.Front()
		conn := p.conn
		p.mutex.Unlock()

		if elem == nil || conn == nil {
			break
		}

		item := elem.Value.(queueItem)

		if item.State {
			err := conn.Subscribe(item.Key)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("cache.pubsub: Subscribe error")
				time.Sleep(constants.RetryDelay)
				continue
			}
		} else {
			err := conn.Unsubscribe(item.Key)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("cache.pubsub: Unsubscribe error")
				time.Sleep(constants.RetryDelay)
				continue
			}
		}

		p.mutex.Lock()
		p.queue.Remove(elem)
		p.mutex.Unlock()
	}
}

func (p *pubSubConn) Listen() {
	go func() {
		for {
			p.parseQueue()

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

			time.Sleep(constants.RetryDelay)

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

			p.mutex.Lock()

			p.conn = &redis.PubSubConn{
				Conn: conn,
			}

			p.queue = list.New()

			for key, _ := range p.listeners {
				p.queue.PushBack(queueItem{
					State: true,
					Key:   key,
				})
			}

			p.mutex.Unlock()

			break
		}
	}
}

func (p *pubSubConn) Close() {
	p.mutex.Lock()
	p.closed = true

	conn := p.conn
	if conn != nil {
		p.conn.Close()
		p.conn = nil
	}

	p.mutex.Unlock()
}

func newPubSubConn(address string) (psc *pubSubConn) {
	psc = &pubSubConn{
		address:   address,
		mutex:     &sync.Mutex{},
		queue:     list.New(),
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
