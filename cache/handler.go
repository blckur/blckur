package cache

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
		h.listener.stream <- evt.Data
	case RESHARD:
		h.listener.reshard(h)
	}
}
