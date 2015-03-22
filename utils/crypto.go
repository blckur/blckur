package utils

import (
	"crypto/rand"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/Sirupsen/logrus"
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

func seedRand() {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("utils: Error seeding random")
		mathrand.Seed(time.Now().UTC().UnixNano())
		return
	}

	mathrand.Seed(n.Int64())

	return
}

func init() {
	seedRand()
}
