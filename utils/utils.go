package utils

import (
	"math/rand"
	"github.com/dropbox/godropbox/errors"
	"net/mail"
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
