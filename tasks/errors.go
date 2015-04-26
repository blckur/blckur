package tasks

import (
	"github.com/dropbox/godropbox/errors"
)

type ApiError struct {
	errors.DropboxError
}

type ParseError struct {
	errors.DropboxError
}
