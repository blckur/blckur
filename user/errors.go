package user

import (
    "github.com/dropbox/godropbox/errors"
)

type NotFoundError struct {
    errors.DropboxError
}

type DatabaseError struct {
    errors.DropboxError
}

type UnknownError struct {
    errors.DropboxError
}
