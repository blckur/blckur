package account

import (
	"github.com/dropbox/godropbox/errors"
)

type InvalidTypeError struct {
	errors.DropboxError
}
