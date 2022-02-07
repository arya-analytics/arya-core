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

type ErrorHandler struct {
	ConverterChain   ErrorTypeConverterChain
	DefaultConverter ErrorTypeConverter
}

func NewErrorHandler(cc ErrorTypeConverterChain, dc ErrorTypeConverter) ErrorHandler {
	return ErrorHandler{cc, dc}
}

func (eh ErrorHandler) Exec(err error) error {
	t := reflect.TypeOf(err)
	if err == nil {
		return nil
	}
	if t == reflect.TypeOf(Error{}) {
		return err
	}
	var (
		errT ErrorType
		ok   bool
	)
	for _, c := range eh.ConverterChain {
		errT, ok = c(err)
		if ok {
			break
		}
	}
	if !ok {
		errT, ok = eh.DefaultConverter(err)
	}

	if !ok {
		return unknownErr(err)
	}
	return Error{
		Type: errT,
		Base: err,
	}
}

func unknownErr(err error) error {
	log.Errorf("Storage - Unknown Err -> %s", err)
	return Error{
		Type: ErrTypeUnknown,
		Base: err,
	}
}
