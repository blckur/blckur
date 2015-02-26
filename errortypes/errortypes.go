package errortypes

import (
    "github.com/dropbox/godropbox/errors"
)

type UnknownError struct {
    errors.DropboxError
}
