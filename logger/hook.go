package logger

import (
	"github.com/Sirupsen/logrus"
)

type logHook struct{}

func (h *logHook) Fire(entry *logrus.Entry) (err error) {
	if len(paperTrailBuffer) <= 125 {
		paperTrailBuffer <- entry
	}

	if len(rollbarBuffer) <= 125 {
		rollbarBuffer <- entry
	}

	return
}

func (h *logHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
