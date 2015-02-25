package auth

import (
    "github.com/dropbox/godropbox/errors"
)

type NotFoundError struct {
    errors.DropboxError
}

type UnknownError struct {
    errors.DropboxError
}
