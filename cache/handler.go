package cache

type handler struct {
	listener *Listener
	State bool
}

func (h *handler) Handle(evt *Event) {
	if !h.State {
		return
	}

	switch evt.Type {
	case MESSAGE:
		h.listener.Stream <- evt.Data
	case RESHARD:
		h.listener.reshard(h)
	}
}
