package utils

import (
	"crypto/rand"
	"github.com/dropbox/godropbox/errors"
	"net/mail"
)

func RandBytes(size int) (bytes []byte, err error) {
	bytes = make([]byte, 32)
	_, err = rand.Read(bytes)
	return
}

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
