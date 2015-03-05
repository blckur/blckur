package requires

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
)

var (
	loaded = set.NewSet()
)

func Before(name string) {
	if loaded.Contains(name) {
		err := &InitError{
			errors.New( "settings: Init before error"),
		}
		panic(err)
	}
	return
}

func After(name string) {
	if !loaded.Contains(name) {
		err := &InitError{
			errors.New("settings: Init after error"),
		}
		panic(err)
	}
	return
}

func Register(name string) {
	loaded.Add(name)
}
