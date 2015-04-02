// Runs defined tasks at scheduled time.
package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/tasks"
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

type Reserve struct {
	Id        uint32    `bson:"_id"`
	Timestamp time.Time `bson:"timestamp"`
}

func runTasks(tme time.Time) {
	tasks := tasks.GetTasks(tme)

	if len(tasks) == 0 {
		return
	}

	db := database.GetDatabase()
	coll := db.Tasks()

	for _, task := range tasks {
		idHash := fnv.New32a()
		idHash.Write([]byte(task.Type()))
		idHash.Write([]byte(strconv.FormatInt(tme.Unix(), 10)))
		id := idHash.Sum32()

		resv := &Reserve{
			Id:        id,
			Timestamp: time.Now(),
		}

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		err := coll.Insert(resv)
		if err != nil {
			err = database.ParseError(err)
			switch err.(type) {
			case *database.DuplicateKeyError:
				err = nil
			default:
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("scheduler: Reserve error")
			}
			return
		}

		err = task.Run(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("scheduler: Task run error")
		}
	}
}

func Run() {
	lastCheck := time.Now()
	for {
		curTime := time.Now()
		curTime = curTime.Add(-time.Duration(
			curTime.Nanosecond()) * time.Nanosecond)
		curTime = curTime.Add(-time.Duration(
			curTime.Second()) * time.Second)

		if curTime != lastCheck {
			lastCheck = curTime
			go runTasks(curTime)
		}

		time.Sleep(1 * time.Second)
	}
}
