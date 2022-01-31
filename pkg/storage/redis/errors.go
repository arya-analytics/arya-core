package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
)

func parseRedisTSErr(err error) (oErr error) {
	if err == nil {
		return err
	}
	switch err.(type) {
	case storage.Error:
		oErr = err
	default:
		oErr = storage.Error{Type: redisToStorageErrType(err), Base: err,
			Message: err.Error()}
	}
	se, ok := oErr.(storage.Error)
	if ok {
		if se.Type == storage.ErrTypeUnknown {
			log.Errorf("Unknown err -> %s", err)
		}
	}
	return oErr
}

func redisToStorageErrType(err error) storage.ErrorType {
	ot, ok := _redisErrors[err.Error()]
	if !ok {
		return storage.ErrTypeUnknown
	}
	return ot
}

var _redisErrors = map[string]storage.ErrorType{
	"ERR TSDB: the key does not exist": storage.ErrTypeItemNotFound,
	"ERR TSDB: key already exists":     storage.ErrTypeUniqueViolation,
}
