package cache

import (
	"container/list"
	"github.com/blckur/blckur/constants"
	"github.com/garyburd/redigo/redis"
	"github.com/dropbox/godropbox/errors"
	"github.com/Sirupsen/logrus"
	"sync"
	"time"
)

const (
	MESSAGE = 1
	RESHARD = 2
)

type Event struct {
	Type int
	Data string
}

type PubSubConn struct {
	address string
	addMutex *sync.Mutex
	remMutex *sync.Mutex
	psConn *redis.PubSubConn
	listeners map[string]func(evt *Event)
	addListeners *list.List
	remListeners *list.List
	closed bool
}

func (p *PubSubConn) subscribe(key string, handler func(evt *Event)) {
	p.addMutex.Lock()
	p.listeners[key] = handler
	p.addListeners.PushBack(key)
	p.addMutex.Unlock()
}

func (p *PubSubConn) unsubsribe(key string) {
	p.remMutex.Lock()
	p.remListeners.PushBack(key)
	delete(p.listeners, key)
	p.remMutex.Unlock()
}

func (p *PubSubConn) reshard() {
	p.addMutex.Lock()
	p.remMutex.Lock()

	evt := &Event{
		Type: RESHARD,
	}

	for _, handler := range p.listeners {
		handler(evt)
	}

	p.close()

	p.addMutex.Unlock()
	p.remMutex.Unlock()
}

func (p *PubSubConn) listen() {
	go func() {
		for {
			if p.closed {
				return
			}

			for {
				p.addMutex.Lock()
				elem := p.addListeners.Front()
				conn := p.psConn
				p.addMutex.Unlock()

				if elem == nil || conn.Conn == nil {
					break
				}

				err := conn.Subscribe(elem.Value.(string))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("cache: Subscribe error")
					time.Sleep(constants.RETRY_DELAY)
					continue
				}

				p.addMutex.Lock()
				p.addListeners.Remove(elem)
				p.addMutex.Unlock()
			}

			for {
				p.remMutex.Lock()
				elem := p.remListeners.Front()
				conn := p.psConn
				p.remMutex.Unlock()

				if elem == nil || conn.Conn == nil {
					break
				}

				err := conn.Unsubscribe(elem.Value.(string))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("cache: Unsubscribe error")
					time.Sleep(constants.RETRY_DELAY)
					continue
				}

				p.remMutex.Lock()
				p.remListeners.Remove(elem)
				p.remMutex.Unlock()
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		if p.psConn.Conn != nil {
			Loop:
			for {
				switch obj := p.psConn.Receive().(type) {
				case redis.Message:
					if listener, ok := p.listeners[obj.Channel]; ok {
						listener(&Event{
							Type: MESSAGE,
							Data: string(obj.Data),
						})
					}
				case error:
					if p.closed {
						return
					}

					logrus.WithFields(logrus.Fields{
						"error": obj,
					}).Error("cache: Listen error")

					break Loop
				}
			}
		}

		if p.closed {
			return
		}

		for {
			time.Sleep(constants.RETRY_DELAY)

			conn, err := dialLong(p.address)
			if err != nil {
				err = &CacheError{
					errors.Wrap(err, "cache: Dial error"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("cache: Dial error")
				continue
			}

			p.addMutex.Lock()
			p.remMutex.Lock()

			p.psConn = &redis.PubSubConn{
				Conn: conn,
			}

			p.addListeners = list.New()
			p.remListeners = list.New()

			for key, _ := range p.listeners {
				p.addListeners.PushBack(key)
			}

			p.addMutex.Unlock()
			p.remMutex.Unlock()

			break
		}
	}
}

func (p *PubSubConn) close() {
	p.closed = true
	p.psConn.Close()
	p.psConn = nil
}

func newPubSubConn(address string) (psc *PubSubConn) {
	psc = &PubSubConn{
		address: address,
		addMutex: &sync.Mutex{},
		remMutex: &sync.Mutex{},
		listeners: map[string]func(evt *Event){},
		addListeners: list.New(),
		remListeners: list.New(),
		closed: false,
	}

	conn, err := dialLong(address)
	if err != nil {
		err = &CacheError{
			errors.Wrap(err, "cache: Connection error"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("cache: Connection error")
	}
	psc.psConn = &redis.PubSubConn{
		Conn: conn,
	}

	return
}
