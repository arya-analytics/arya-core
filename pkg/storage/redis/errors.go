package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
)

var _errTypeConverterChain = storage.ErrorTypeConverterChain{}
var _defaultConverter = errConverterDefault

func errConverterDefault(err error) (storage.ErrorType, bool) {
	ot, ok := _redisErrors[err.Error()]
	return ot, ok
}

var _redisErrors = map[string]storage.ErrorType{
	"ERR TSDB: the key does not exist": storage.ErrTypeItemNotFound,
	"ERR TSDB: key already exists":     storage.ErrTypeUniqueViolation,
}

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(_errTypeConverterChain, _defaultConverter)
}
