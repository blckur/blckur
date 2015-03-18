package nodes

import (
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/utils"
	"os/exec"
	"time"
	"strconv"
)

type BeanstalkdNode struct {
	Id string
	Host string
	Port int
}

func (b *BeanstalkdNode) Start() {
	portInt := b.Port
	if portInt == 0 {
		portInt = 11300
	}
	port := strconv.Itoa(portInt)

	args := []string{"-p", port}
	if b.Host != "" {
		args = append(args, "-l", b.Host)
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
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			cmd.Wait()
			time.Sleep(constants.RETRY_DELAY)
		}

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
				}).Error("beanstalkd: Unexpected exit")
				break
			}

			_, err := coll.UpsertId(b.Id, &Node{
				Id: b.Id,
				Type: "beanstalkd",
				Address: address,
				Timestamp: time.Now(),
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
