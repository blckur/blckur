package nodes

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/node"
	"github.com/Sirupsen/logrus"
	"os/exec"
	"time"
	"strconv"
)

type QueueNode struct {
	Id string
	Host string
	Port int
}

func (q *QueueNode) Start() {
	port := strconv.Itoa(q.Port)

	args := []string{"-p", port}
	if q.Host != "" {
		args = append(args, "-l", q.Host)
	}

	address, err := utils.GetLocalAddress()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("queue: Failed to get ip")
		panic(err)
	}
	address = address + ":" + port

	var cmd *exec.Cmd
	for {
		if cmd != nil {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			cmd.Wait()
			time.Sleep(constants.RETRY_DELAY)
		}

		logrus.WithFields(logrus.Fields{
			"id": q.Id,
			"address": address,
		}).Info("nodes.queue: Starting queue node")

		cmd = exec.Command("./bin/beanstalkd", args...) // TODO
		cmd.Stdout = logger.NewErrorWriter()
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
				time.Sleep(constants.RETRY_DELAY)
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
				Id: q.Id,
				Type: "queue",
				Address: address,
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
