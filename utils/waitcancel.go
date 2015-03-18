package utils

type WaitCancel struct {
	count int
	wait chan bool
}

func (w *WaitCancel) Add(n int) {
	if w.wait == nil {
		w.wait = make(chan bool)
	}
	w.count += n
}

func (w *WaitCancel) Done() {
	w.wait <- false
}

func (w *WaitCancel) Cancel() {
	w.wait <- true
}

func (w *WaitCancel) Wait() {
	done := 0
	for {
		if <-w.wait {
			return
		} else {
			done += 1
			if done >= w.count {
				return
			}
		}
	}
}
