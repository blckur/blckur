package cache

import (
	"github.com/dropbox/godropbox/errors"
)

type CacheError struct {
	errors.DropboxError
}
