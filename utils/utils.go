package utils

import (
	"math/rand"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"net"
	"net/mail"
	"os"
)

func ParseEmail(input string) (email string, err error) {
	address, err := mail.ParseAddress(input)
	if err != nil {
		err = &InvalidEmailError{
			errors.Wrap(err, "utils: Invalid email"),
		}
		return
	}
	email = address.Address
	return
}

func ShuffleStrings(data []string) {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func ShuffleStringsNew(data []string) (randData []string) {
	randData = make([]string, len(data))
	copy(randData, data)
	ShuffleStrings(randData)
	return
}

func GetLocalAddress() (addr string, err error) {
	name, err := os.Hostname()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "beanstalkd: Get ip"),
		}
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "beanstalkd: Get ip"),
		}
		return
	}

	addr = addrs[0]

	return
}
