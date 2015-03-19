package cache

import (
	"sync"
)

type Listener struct {
	mutex *sync.Mutex
	clst *cluster
	key string
	servers []string
	handlers []*handler
	Stream chan string
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

func (l *Listener) Close() {
	l.mutex.Lock()
	for _, handler := range l.handlers {
		handler.State = false
	}

	for _, server := range l.servers {
		l.clst.pubsubConns[server].unsubsribe(l.key)
	}

	close(l.Stream)
	l.mutex.Unlock()
}

func (l *Listener) sub() {
	cst := clst
	servers := []string{}
	handlers := []*handler{}

	for _, server := range cst.shrd.Select(l.key) {
		servers = append(servers, server)

		handler := &handler{
			listener: l,
			State: true,
		}
		handlers = append(handlers, handler)

		cst.pubsubConns[server].subscribe(l.key, handler.Handle)
	}

	l.clst = cst
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

func Subscribe(key string) (lst *Listener) {
	lst = &Listener{
		mutex: &sync.Mutex{},
		key: key,
		Stream: make(chan string),
	}

	lst.init()

	return
}
