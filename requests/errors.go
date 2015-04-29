package requests

import (
	"github.com/dropbox/godropbox/errors"
)

type ParseError struct {
	errors.DropboxError
}

type RequestError struct {
	errors.DropboxError
}

type ResponseError struct {
	errors.DropboxError
}
