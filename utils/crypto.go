package utils

import (
	"crypto/rand"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/logger"
	"github.com/dropbox/godropbox/errors"
	"math/big"
	"time"
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

func SeedRand() {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		logger.Error("Failed to seed random with urandom: %s", err)
		mathrand.Seed(time.Now().UTC().UnixNano())
		return
	}

	mathrand.Seed(n.Int64())

	return
}
