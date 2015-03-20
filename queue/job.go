package queue

import (
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
)

type Job struct {
	beanId uint64
	conn *beanstalk.Conn
	Id bson.ObjectId
	Type string
	Data interface{}
}

func (j *Job) Delete() (err error) {
	err = j.conn.Delete(j.beanId)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue: Unknown error"),
		}
		return
	}

	return
}
