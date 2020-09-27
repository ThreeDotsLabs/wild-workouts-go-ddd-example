package logs

import (
	"github.com/sirupsen/logrus"
)

func LogCommandExecution(commandName string, cmd interface{}, err error) {
	log := logrus.WithField("cmd", cmd)

	if err == nil {
		log.Info(commandName + " command succeeded")
	} else {
		log.WithError(err).Error(commandName + " command failed")
	}
}
