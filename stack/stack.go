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
	defer s.mutex.Unlock()

	n := len(s.data)
	if n == 0 {
		return
	}

	elem = s.data[n - 1]
	s.data = s.data[:n - 1]

	return
}

func NewStringStack(data []string) *StringStack {
	return &StringStack{
		data: data,
		mutex: &sync.Mutex{},
	}
}
