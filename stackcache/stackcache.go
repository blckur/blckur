// Caches elements FIFO in set with limit.
package stackcache

import (
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"sync"
)

type IdStackCache struct {
	queue chan bson.ObjectId
	data set.Set
	limit int
	mutex *sync.Mutex
}

func (s *IdStackCache) Add(elem bson.ObjectId) {
	s.mutex.Lock()

	if len(s.queue) >= s.limit {
		s.data.Remove(<-s.queue)
	}

	s.data.Add(elem)
	s.queue <- elem

	s.mutex.Unlock()
}

func (s *IdStackCache) Contains(elem bson.ObjectId) bool {
	return s.data.Contains(elem)
}

func NewIdStackCache(limit int) (stck *IdStackCache) {
	stck = &IdStackCache{
		queue: make(chan bson.ObjectId, limit),
		data: set.NewSet(),
		limit: limit,
		mutex: &sync.Mutex{},
	}
	return
}
