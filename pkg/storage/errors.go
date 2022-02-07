package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
)

// |||| ERROR TYPES ||||

const (
	errKey = "storage"
)

type Error struct {
	Base    error
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s - %s", errKey, e.Type, e.Message)
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeItemNotFound
	ErrorTypeUniqueViolation
	ErrorTypeRelationshipViolation
	ErrorTypeInvalidField
	ErrorTypeNoPK
	ErrorTypeMigration
	ErrorTypeInvalidArgs
)

type ErrorTypeConverter func(err error) (ErrorType, bool)

func newIterConverter(ecc []ErrorTypeConverter) func() (ErrorTypeConverter, bool) {
	n := -1
	return func() (ErrorTypeConverter, bool) {
		n += 1
		if n < len(ecc) {
			return ecc[n], true
		}
		return nil, false
	}
}

type ErrorHandler struct {
	ConverterChain   []ErrorTypeConverter
	DefaultConverter ErrorTypeConverter
}

func NewErrorHandler(dc ErrorTypeConverter, cc ...ErrorTypeConverter) ErrorHandler {
	return ErrorHandler{cc, dc}
}

func (eh ErrorHandler) Exec(err error) error {
	if err == nil || isStorageError(err) {
		return err
	}
	errT, ok := eh.ErrorType(err)
	if !ok {
		return unknownErr(err)
	}
	return Error{
		Type: errT,
		Base: err,
	}
}

func (eh ErrorHandler) ErrorType(err error) (ErrorType, bool) {
	next := newIterConverter(eh.ConverterChain)
	for {
		c, nextOk := next()
		if !nextOk {
			return eh.DefaultConverter(err)
		}
		if errT, ok := c(err); ok {
			return errT, ok
		}
	}
}

func unknownErr(err error) error {
	log.Errorf("Storage - Unknown Err -> %s", err)
	return Error{
		Type:    ErrorTypeUnknown,
		Base:    err,
		Message: "storage - unknown error",
	}
}

func isStorageError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(Error{})
}
