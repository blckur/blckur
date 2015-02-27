package auth

import (
	"github.com/dropbox/godropbox/errors"
)

type NotFoundError struct {
	errors.DropboxError
}
