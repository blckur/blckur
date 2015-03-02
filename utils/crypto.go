package utils

import (
	"crypto/rand"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"math/big"
	mathrand "math/rand"
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

func SeedRand() (err error) {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "utils: Random seed error"),
		}
		return
	}

	mathrand.Seed(n.Int64())

	return
}
