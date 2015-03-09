package database

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"fmt"
)

func GetErrorCode(err error) (errCode int) {
	switch err := err.(type) {
	case *mgo.LastError:
		errCode = err.Code
	case *mgo.QueryError:
		errCode = err.Code
	}

	return
}

func ParseError(err error) (newErr error) {
	if err == mgo.ErrNotFound {
		newErr = &NotFoundError{
			errors.New("database: Not found"),
		}
		return
	}

	errCode := GetErrorCode(err)

	switch errCode {
	case 11000, 11001, 12582, 16460:
		newErr = &DuplicateKeyError{
			errors.New("database: Duplicate key"),
		}
	default:
		newErr = &UnknownError{
			errors.Wrap(err, fmt.Sprintf(
				"database: Unknown error %d", errCode)),
		}
	}

	return
}

func ExpectFound(c *gin.Context, err error) {
	switch err.(type) {
	case NotFoundError:
		c.Fail(404, err)
	default:
		c.Fail(500, err)
	}
}
