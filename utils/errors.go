package utils

import (
	"github.com/dropbox/godropbox/errors"
)

type InvalidEmailError struct {
	errors.DropboxError
}

type InitError struct {
	errors.DropboxError
}
