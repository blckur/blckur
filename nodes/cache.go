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

type CacheNode struct {
	Id   string
	Host string
	Port int
}

func (c *CacheNode) Start() {
	port := strconv.Itoa(c.Port)
	address := getAddress() + ":" + port

	args := []string{
		"--save", "",
		"--port", port,
	}

	if c.Host != "" {
		args = append(args, "--host", c.Host)
	}

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
			"id":      c.Id,
			"address": address,
		}).Info("nodes.cache: Starting cache node")

		cmd = exec.Command("redis-server", args...)
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
				}).Error("cache: Unexpected exit")
				break
			}

			stat, err := coll.UpsertId(c.Id, &node.Node{
				Id:        c.Id,
				Type:      "cache",
				Address:   address,
				Timestamp: time.Now(),
			})
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("cache: Database upsert")
				continue
			}

			time.Sleep(5 * time.Second)

			if stat.Updated == 0 {
				messenger.Publish(db, "cache", "update")
			}
		}
	}
}
