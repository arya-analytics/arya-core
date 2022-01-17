package validate

import (
	"fmt"
	"reflect"
)

type ValidateFunc func(v reflect.Value) error

type Validator struct {
	validators []ValidateFunc
}

type Error struct {
	Type ErrorType
}

func (err Error) Error() string {
	return fmt.Sprint(err.Type)
}

type ErrorType int

func NewError(t ErrorType) Error {
	return Error{t}
}

const (
	NonPointer = iota
	NonStructOrSlice
)

func (v *Validator) Exec(m interface{}) error {
	for _, v := range v.validators {
		if err := v(reflect.ValueOf(m)); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSliceOrStruct(v reflect.Value) error {
	k := v.Type().Elem().Kind()
	if k != reflect.Struct && k != reflect.Slice {
		return NewError(NonStructOrSlice)
	}
	return nil
}

func ValidateContainerIsPointer(v reflect.Value) error {
	k := v.Type().Kind()
	if k != reflect.Ptr {
		return NewError(NonPointer)
	}
	return nil
}
