package utils

import (
	"crypto/rand"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
)

func RandBytes(size int) (bytes []byte, err error) {
	bytes = make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "utils: Random read error"),
		}
		return
	}

	return
}
