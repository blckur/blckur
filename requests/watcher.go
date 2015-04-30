package requests

import (
	"bufio"
	"net/http"
	"time"
)

type Event struct {
	Type string
	Data string
}

type Watcher struct {
	scanner *bufio.Scanner
	res     *http.Response
	timeout time.Duration
	Stream  chan *Event
	Err     error
}

func (w *Watcher) scan() {
	defer close(w.Stream)
	defer w.Close()

	var typ string
	timer := time.AfterFunc(w.timeout, func() {
		w.Close()
	})

	for w.scanner.Scan() {
		data := w.scanner.Text()
		n := len(data)

		if n >= 7 && data[:7] == "event: " {
			typ = data[7:]
			continue
		} else if n >= 6 && data[:6] == "data: " {
			data = data[6:]
		} else {
			continue
		}

		if typ == "keep-alive" {
			timer.Reset(w.timeout)
			continue
		}

		w.Stream <- &Event{
			Type: typ,
			Data: data,
		}
	}
}

func (w *Watcher) Close() {
	w.res.Body.Close()
}
