package nodes

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/messenger"
	"os/exec"
	"time"
	"strconv"
)

type CacheNode struct {
	Id string
	Host string
	Port int
}

func (c *CacheNode) Start() {
	portInt := c.Port
	if portInt == 0 {
		portInt = 6379
	}
	port := strconv.Itoa(portInt)

	args := []string{
		"--save", "",
		"--port", port,
	}

	if c.Host != "" {
		args = append(args, "--host", c.Host)
	}

	address, err := utils.GetLocalAddress()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("cache: Failed to get ip")
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

		cmd = exec.Command("redis-server", args...)
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
				}).Error("cache: Unexpected exit")
				break
			}

			stat, err := coll.UpsertId(c.Id, &Node{
				Id: c.Id,
				Type: "cache",
				Address: address,
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
