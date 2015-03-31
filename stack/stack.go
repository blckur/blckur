// Simple stack with push and pop.
package stack

import (
	"sync"
)

type StringStack struct {
	data []string
	mutex *sync.Mutex
}

func (s *StringStack) Push(elem string) {
	s.mutex.Lock()
	s.data = append(s.data, elem)
	s.mutex.Unlock()
}

func (s *StringStack) Pop() (elem string) {
	s.mutex.Lock()

	n := len(s.data)
	if n == 0 {
		return
	}

	elem = s.data[n - 1]
	s.data = s.data[:n - 1]

	s.mutex.Unlock()
	return
}

func NewStringStack(data []string) (stck *StringStack) {
	stck = &StringStack{
		data: data,
		mutex: &sync.Mutex{},
	}
	return
}
