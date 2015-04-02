package queue

import (
	"github.com/dropbox/godropbox/container/set"
)

type Listener struct {
	servers  set.Set
	streams  map[string]*Stream
	onStream func(*Stream)
}

func (l *Listener) updateStreams(servers set.Set) {
	add := servers.Copy()
	add.Subtract(l.servers)
	rem := l.servers.Copy()
	rem.Subtract(servers)

	for serverInf := range add.Iter() {
		server := serverInf.(string)
		stream := newStream(server)

		go l.onStream(stream)
		l.streams[server] = stream
	}

	for serverInf := range rem.Iter() {
		server := serverInf.(string)
		l.streams[server].Stop = true
		delete(l.streams, server)
	}

	l.servers = servers
}

func NewListener(onStream func(*Stream)) (lstnr *Listener) {
	lstnr = &Listener{
		servers:  set.NewSet(),
		streams:  map[string]*Stream{},
		onStream: onStream,
	}

	mutex.Lock()
	listeners = append(listeners, lstnr)
	lstnr.updateStreams(clst.servers)
	mutex.Unlock()

	return
}
