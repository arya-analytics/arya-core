package errutil

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

func Inspect(err error) (error, bool) {
	log.WithFields(log.Fields{
		"type": reflect.TypeOf(err),
	}).Infof("Inspection of error -> %s", err)
	return err, false
}
