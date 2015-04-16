package email

import (
	"github.com/dropbox/godropbox/errors"
)

type TlsError struct {
	errors.DropboxError
}

type SmtpError struct {
	errors.DropboxError
}
