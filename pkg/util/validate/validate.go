package validate

import (
	"reflect"
)

type ValidateFunc func(v reflect.Value) error

type Validator struct {
	validators []ValidateFunc
}

func New(v []ValidateFunc) *Validator {
	return &Validator{v}
}

func (v *Validator) Exec(m interface{}) error {
	for _, v := range v.validators {
		if err := v(reflect.ValueOf(m)); err != nil {
			return err
		}
	}
	return nil
}
