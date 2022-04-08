package auth

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Error struct {
	Type    ErrorType
	Message string
	Base    error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s - %s", e.Type, e.Message)
}

func newSimpleError(errType ErrorType, base error) error {
	return Error{Type: errType, Message: base.Error(), Base: base}
}

//go:generate stringer -type=ErrorType
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeUserNotFound
	ErrorTypeInvalidCredentials
)

func newErrorConvert() errutil.ConvertChain {
	return errutil.ConvertChain{errorConvertQuery, errorConvertDefault}
}

func queryErrors() map[query.ErrorType]ErrorType {
	return map[query.ErrorType]ErrorType{
		query.ErrorTypeItemNotFound: ErrorTypeUserNotFound,
	}
}

func errorConvertQuery(err error) (error, bool) {
	qe, ok := err.(query.Error)
	if !ok {
		return err, false
	}
	t, ok := queryErrors()[qe.Type]
	return newSimpleError(t, qe), ok
}

func errorConvertDefault(err error) (error, bool) {
	_, ok := err.(Error)
	if ok {
		return err, true
	}
	return newSimpleError(ErrorTypeUnknown, err), true
}
