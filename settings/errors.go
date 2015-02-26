package settings

import (
    "github.com/dropbox/godropbox/errors"
)

type DatabaseError struct {
    errors.DropboxError
}

type UnknownError struct {
    errors.DropboxError
}
