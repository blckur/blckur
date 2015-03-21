package queue

import (
	"github.com/dropbox/godropbox/container/set"
)

type Listener struct {
	streams set.Set
	onStream func(*Stream)
}

func (l *Listener) setStreams(streams set.Set) {
	add := l.streams.Copy()
	add.Subtract(streams)
	rem := l.streams.Copy()
	rem.Intersect(streams)

	for stream := range add.Iter() {
		l.onStream(stream.(*Stream))
	}

	for stream := range rem.Iter() {
		stream.(*Stream).Stop = true
	}

	l.streams = streams
}

func NewListener(onStream func(*Stream)) (lstnr *Listener) {
	lstnr = &Listener{
		streams: set.NewSet(),
		onStream: onStream,
	}

	mutex.Lock()
	listeners = append(listeners, lstnr)
	mutex.Unlock()

	return
}
