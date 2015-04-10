package cache

import (
	"github.com/blckur/blckur/stackcache"
	"sync"
)

type Listener struct {
	mutex    *sync.Mutex
	cluster  *cluster
	channel  string
	servers  []string
	handlers []*handler
	stream   chan string
	idCache  *stackcache.IdStackCache
}

func (l *Listener) reshard() {
	l.mutex.Lock()

	for _, handler := range l.handlers {
		handler.State = false
	}

	l.sub()
	l.mutex.Unlock()
}

func (l *Listener) Listen() (stream chan string) {
	return l.stream
}

func (l *Listener) Close() {
	clstMutex.RLock()
	subs.Remove(l)
	clstMutex.RUnlock()

	l.mutex.Lock()

	for _, handler := range l.handlers {
		handler.State = false

		l.cluster.pubsubConns[handler.Server].Unsubsribe(l.channel, handler.Id)
	}
	close(l.stream)

	l.mutex.Unlock()
}

func (l *Listener) sub() {
	cst := clst
	servers := cst.shrd.Select(l.channel)
	handlers := []*handler{}

	for _, server := range servers {
		handler := &handler{
			listener: l,
			State:    true,
			Server:   server,
		}
		handlers = append(handlers, handler)

		handler.Id = cst.pubsubConns[server].Subscribe(
			l.channel, handler.Handle)
	}

	l.cluster = cst
	l.servers = servers
	l.handlers = handlers
}

func (l *Listener) init() {
	clstMutex.RLock()
	l.mutex.Lock()
	l.sub()
	subs.Add(l)
	l.mutex.Unlock()
	clstMutex.RUnlock()
}

// Subscribe to a channel connection will auto repair as nodes are added/lost
func Subscribe(channel string) (lst *Listener) {
	lst = &Listener{
		mutex:   &sync.Mutex{},
		channel: channel,
		stream:  make(chan string),
		idCache: stackcache.NewIdStackCache(64),
	}

	lst.init()

	return
}
