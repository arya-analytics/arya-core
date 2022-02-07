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
	ErrTypeUnknown ErrorType = iota
	ErrTypeItemNotFound
	ErrTypeUniqueViolation
	ErrTypeRelationshipViolation
	ErrTypeInvalidField
	ErrTypeNoPK
	ErrTypeMigration
	ErrTypeInvalidArgs
)

type ErrorTypeConverter func(err error) (ErrorType, bool)

type ErrorTypeConverterChain []ErrorTypeConverter

func newIterConverter(ecc ErrorTypeConverterChain) func() (ErrorTypeConverter, bool) {
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
	ConverterChain   ErrorTypeConverterChain
	DefaultConverter ErrorTypeConverter
}

func NewErrorHandler(cc ErrorTypeConverterChain, dc ErrorTypeConverter) ErrorHandler {
	return ErrorHandler{cc, dc}
}

func (eh ErrorHandler) Exec(err error) error {
	if err == nil || isStorageError(err) {
		return err
	}
	errT, ok := eh.errType(err)
	if !ok {
		return unknownErr(err)
	}
	return Error{
		Type: errT,
		Base: err,
	}
}

func (eh ErrorHandler) errType(err error) (ErrorType, bool) {
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
		Type:    ErrTypeUnknown,
		Base:    err,
		Message: "storage - unknown error",
	}
}

func isStorageError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(Error{})
}
