package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
)

var _ErrorTypeConverterChain = storage.ErrorTypeConverterChain{}
var _defaultConverter = errConverterDefault

func errConverterDefault(err error) (storage.ErrorType, bool) {
	ot, ok := _redisErrors[err.Error()]
	return ot, ok
}

var _redisErrors = map[string]storage.ErrorType{
	"ERR TSDB: the key does not exist": storage.ErrorTypeItemNotFound,
	"ERR TSDB: key already exists":     storage.ErrorTypeUniqueViolation,
}

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(_ErrorTypeConverterChain, _defaultConverter)
}
