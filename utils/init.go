package utils

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
)

var (
	loaded = set.NewSet()
)

func Before(name string) (err error) {
	if loaded.Contains(name) {
		err = &InitError{
			errors.Wrap(err, "settings: Init before error"),
		}
	}
	return
}

func After(name string) (err error) {
	if !loaded.Contains(name) {
		err = &InitError{
			errors.Wrap(err, "settings: Init after error"),
		}
	}
	return
}

func Register(name string) {
	loaded.Add(name)
}
