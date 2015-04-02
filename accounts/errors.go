package accounts

import (
	"github.com/dropbox/godropbox/errors"
)

type InvalidDate struct {
	errors.DropboxError
}
