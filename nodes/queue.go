package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/node"
	"os/exec"
	"strconv"
	"time"
)

type QueueNode struct {
	Id         string
	Host       string
	Address    string
	Port       int
	PublicPort int
}

func (q *QueueNode) Start() {
	port := strconv.Itoa(q.Port)
	address := q.Address + ":" + strconv.Itoa(q.PublicPort)

	args := []string{"-p", port}
	if q.Host != "" {
		args = append(args, "-l", q.Host)
	}

	var cmd *exec.Cmd
	for {
		if cmd != nil {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			cmd.Wait()
			time.Sleep(constants.RetryDelay)
		}

		logrus.WithFields(logrus.Fields{
			"id":          q.Id,
			"address":     q.Address,
			"port":        q.Port,
			"public_port": q.Port,
		}).Info("nodes.queue: Starting queue node")

		cmd = exec.Command("beanstalkd", args...)
		cmd.Stderr = logger.NewErrorWriter()

		var cmdErr error
		go func() {
			cmdErr = cmd.Run()
		}()

		db := database.GetDatabase()
		coll := db.Nodes()

		delay := false
		for {
			if delay {
				time.Sleep(constants.RetryDelay)
			} else {
				delay = true
			}

			if cmdErr != nil {
				logrus.WithFields(logrus.Fields{
					"error": cmd.ProcessState.String(),
				}).Error("queue: Unexpected exit")
				break
			}

			stat, err := coll.UpsertId(q.Id, &node.Node{
				Id:        q.Id,
				Type:      "queue",
				Address:   address,
				Timestamp: time.Now(),
			})
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("queue: Database upsert")
				continue
			}

			time.Sleep(5 * time.Second)

			if stat.Updated == 0 {
				messenger.Publish(db, "queue", "update")
			}
		}
	}
}
