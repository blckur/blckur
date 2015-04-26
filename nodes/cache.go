package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/node"
	"github.com/blckur/blckur/utils"
	"os/exec"
	"strconv"
	"time"
)

type CacheNode struct {
	Id         string
	Host       string
	Address    string
	Port       int
	PublicPort int
}

func (c *CacheNode) ping() {
	address := c.Address + ":" + strconv.Itoa(c.PublicPort)
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Nodes()

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

	if stat.Updated == 0 {
		messenger.Publish(db, "cache", "update")
	}
}

func (c *CacheNode) Start() {
	args := []string{
		"--save", "",
		"--port", strconv.Itoa(c.Port),
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
			time.Sleep(constants.RetryDelay)
		}

		logrus.WithFields(logrus.Fields{
			"id":          c.Id,
			"address":     c.Address,
			"port":        c.Port,
			"public_port": c.Port,
			"version":     utils.GetVersion(),
		}).Info("nodes.cache: Starting cache node")

		cmd = exec.Command("redis-server", args...)
		cmd.Stderr = logger.NewErrorWriter()

		var cmdErr error
		go func() {
			cmdErr = cmd.Run()
		}()
		defer cmd.Process.Kill()
		time.Sleep(5 * time.Second)

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
				}).Error("cache: Unexpected exit")
				break
			}

			c.ping()
		}
	}
}
