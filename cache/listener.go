package cache

import (
	"sync"
)

type Listener struct {
	mutex *sync.Mutex
	cluster *cluster
	channel string
	servers []string
	handlers []*handler
	stream chan string
}

func (l *Listener) reshard(hnd *handler) {
	l.mutex.Lock()
	if !hnd.State {
		l.mutex.Unlock()
		return
	}

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
	l.mutex.Lock()
	for _, handler := range l.handlers {
		handler.State = false
	}

	for _, server := range l.servers {
		l.cluster.pubsubConns[server].Unsubsribe(l.channel)
	}

	close(l.stream)
	l.mutex.Unlock()
}

func (l *Listener) sub() {
	cst := clst
	servers := []string{}
	handlers := []*handler{}

	for _, server := range cst.shrd.Select(l.channel) {
		servers = append(servers, server)

		handler := &handler{
			listener: l,
			State: true,
		}
		handlers = append(handlers, handler)

		cst.pubsubConns[server].Subscribe(l.channel, handler.Handle)
	}

	l.cluster = cst
	l.servers = servers
	l.handlers = handlers
}

func (l *Listener) init() {
	clstMutex.RLock()
	l.mutex.Lock()
	l.sub()
	l.mutex.Unlock()
	clstMutex.RUnlock()
}

func Subscribe(channel string) (lst *Listener) {
	lst = &Listener{
		mutex: &sync.Mutex{},
		channel: channel,
		stream: make(chan string),
	}

	lst.init()

	return
}
