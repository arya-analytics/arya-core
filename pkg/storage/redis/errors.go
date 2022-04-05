package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

func newErrorConvert() errutil.ConvertChain {
	return query.NewErrorConvertChain(internal.ErrorConvertConnection, errorConvertDefault)
}

func errorConvertDefault(err error) (error, bool) {
	t, ok := redisErrors()[err.Error()]
	return query.NewSimpleError(t, err), ok
}

func redisErrors() map[string]query.ErrorType {
	return map[string]query.ErrorType{
		"ERR TSDB: the key does not exist": query.ErrorTypeItemNotFound,
		"ERR TSDB: key already exists":     query.ErrorTypeUniqueViolation,
	}
}
