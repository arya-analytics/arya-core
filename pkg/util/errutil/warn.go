package errutil

import log "github.com/sirupsen/logrus"

func Warn(err error) {
	if err != nil {
		log.Warn(err)
	}
}
