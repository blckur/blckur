package logger

import (
	"github.com/Sirupsen/logrus"
	"time"
)

type limiter map[string]time.Time

func (l limiter) Check(entry *logrus.Entry, limit time.Duration) bool {
	// TODO Use hash for key
	if timestamp, ok := l[entry.Message]; ok &&
		time.Since(timestamp) < limit {

		return false
	}
	l[entry.Message] = time.Now()

	return true
}
