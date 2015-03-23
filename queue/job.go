package queue

import (
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/cache"
	"github.com/dropbox/godropbox/errors"
	"github.com/kr/beanstalk"
	"labix.org/v2/mgo/bson"
	"time"
)

type Job struct {
	beanId uint64
	conn *beanstalk.Conn
	cacheConn *cache.ClusterConn
	Id bson.ObjectId `json:"id"`
	Type string `json:"type"`
	Ttl time.Duration `json:"ttl"`
	Timestamp time.Time `json:"timestamp"`
	Resource bson.ObjectId `json:"resource"`
	Data map[string]string `json:"data"`
}

func (j *Job) Put(priority int, delay time.Duration,
		ttr time.Duration) (err error) {
	err = put(j, priority, delay, ttr)
	return
}

func (j *Job) Delete() (err error) {
	if j.conn == nil {
		return
	}

	err = j.conn.Delete(j.beanId)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "queue.job: Job delete error"),
		}
		return
	}

	err = j.cacheConn.SetString(j.Id.Hex(), "t", 3 * time.Minute)
	if err != nil {
		return
	}

	return
}

func NewJob(ttl time.Duration) (job *Job) {
	job = &Job{
		Id: bson.NewObjectId(),
		Ttl: ttl,
		Timestamp: time.Now(),
	}
	return
}
