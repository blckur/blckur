// Starting point for node types.
package nodes

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/utils"
)

func getAddress() (address string) {
	address, err := utils.GetLocalAddress()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("queue: Failed to get ip")
		panic(err)
	}
	return
}
