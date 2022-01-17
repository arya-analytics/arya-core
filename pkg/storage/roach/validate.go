package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
)

const (
	pkFieldName = "ID"
)

func validatePK(v reflect.Value) (err error) {
	if storage.IsChainModel(v.Elem().Type()) {
		for i := 0; i < v.Elem().Len(); i++ {
			err = validatePK(v.Elem().Index(i))
		}
	} else {
		f := v.Elem().FieldByName(pkFieldName)
		switch f.Kind() {
		case reflect.Int:
			if f.Interface() == 0 {
				err = storage.NewError(storage.ErrTypeNoPK)
			}
		}
	}
	return err
}

var createValidator = validate.New([]validate.ValidateFunc{
	validatePK,
})
