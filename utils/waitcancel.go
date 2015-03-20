package utils

import (
	"sync/atomic"
)

// Warning: Add() can not be called after Wait()
type WaitCancel struct {
	count int32
	wait chan int
}

func (w *WaitCancel) Add(n int) {
	if w.wait == nil {
		w.wait = make(chan int)
	}
	atomic.AddInt32(&w.count, int32(n))
}

func (w *WaitCancel) Done() {
	val := atomic.AddInt32(&w.count, -1)
	if val == -1 {
		w.wait <- 1
	}
}

func (w *WaitCancel) Cancel() {
	val := atomic.SwapInt32(&w.count, -1073741824)
	if val >= 0 {
		w.wait <- 1
	}
}

func (w *WaitCancel) Wait() {
	val := atomic.AddInt32(&w.count, -1)
	if val == -1 {
		return
	}
	<-w.wait
}
