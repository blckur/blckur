package cache

import (
	"encoding/json"
)

type handler struct {
	Id       int
	State    bool
	Server   string
	listener *Listener
}

func (h *handler) Handle(data string) {
	if !h.State {
		return
	}

	defer func() {
		recover()
	}()

	msg := &Message{}
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		panic(err)
		return
	}

	if h.listener.idCache.Contains(msg.Id) {
		return
	}
	h.listener.idCache.Add(msg.Id)

	h.listener.stream <- msg
}
