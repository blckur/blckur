package accounts

import (
	"github.com/dropbox/godropbox/errors"
)

type ParseError struct {
	errors.DropboxError
}

type InvalidDate struct {
	errors.DropboxError
}
