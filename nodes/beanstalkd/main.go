package main

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/logger"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/utils"
	"labix.org/v2/mgo/bson"
	"os/exec"
	"time"
	"os"
)

func main() {
	database.Init()
	settings.Init()
	logger.Init()

	id := bson.NewObjectId()
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "11300"
	}

	args := []string{"-p", port}
	if host != "" {
		args = append(args, "-l", host)
	}

	address, err := utils.GetLocalAddress()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("beanstalkd: Failed to get ip")
		panic(err)
	}
	address = address + ":" + port

	var cmd *exec.Cmd
	for {
		if cmd != nil {
			cmd.Process.Kill()
			cmd.Wait()
			time.Sleep(constants.RETRY_DELAY)
		}

		cmd = exec.Command("./beanstalkd2", args...) // TODO
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
				}).Error("beanstalkd: Unexpected exit")
				break
			}

			_, err := coll.UpsertId(id, &bson.M{
				"type": "beanstalkd",
				"address": address,
				"timestamp": time.Now(),
			})
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("beanstalkd: Database upsert")
				continue
			}

			time.Sleep(5 * time.Second)
		}
	}
}
