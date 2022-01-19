package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
)

const (
	pkFieldName = "ID"
)

func validatePK(v interface{}) (err error) {
	rfl := v.(*model.Reflect)
	if rfl.IsChain() {
		for i := 0; i < rfl.ChainValue().Len(); i++ {
			err = validatePK(rfl.ChainValueByIndex(i))
		}
	} else {
		f := rfl.Value().FieldByName(pkFieldName)
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
