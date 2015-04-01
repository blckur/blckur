package cache

import (
	"encoding/json"
)

type handler struct {
	State bool
	listener *Listener
}

func (h *handler) Handle(evt *event) {
	if !h.State {
		return
	}

	switch evt.Type {
	case MESSAGE:
		defer func() {
			recover()
		}()

		msg := &message{}
		err := json.Unmarshal([]byte(evt.Data), msg)
		if err != nil {
			panic(err)
			return
		}

		if h.listener.idCache.Contains(msg.Id) {
			return
		}
		h.listener.idCache.Add(msg.Id)

		h.listener.stream <- msg.Data
	case RESHARD:
		h.listener.reshard(h)
	}
}
