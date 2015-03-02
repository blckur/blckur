package queue

import (
	"github.com/dropbox/godropbox/errors"
)

type JobFailed struct {
	errors.DropboxError
}
