package internal

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"net"
	"net/url"
)

func ErrorConvertConnection(err error) (error, bool) {
	switch pErr := err.(type) {
	case *net.OpError:
		return query.NewSimpleError(query.ErrorTypeConnection, pErr), true
	case *url.Error:
		return query.NewSimpleError(query.ErrorTypeConnection, pErr), true
	default:
		return pErr, false
	}
}
