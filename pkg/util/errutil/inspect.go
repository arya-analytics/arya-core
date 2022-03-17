package errutil

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

// Inspect logs the following information about an error:
//		1. Its type :)
// Adding more soon.
func Inspect(err error) (error, bool) {
	log.WithFields(log.Fields{
		"type": reflect.TypeOf(err),
	}).Infof("Inspection of error -> %s", err)
	return err, false
}
