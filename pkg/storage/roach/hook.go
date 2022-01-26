package roach

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
)

// || BINDING ||

func beforeInsertSetUUID(rfl *model.Reflect) *model.Reflect {
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		fldT, _ := nRfl.Type().FieldByName(model.KeyPK)
		fld := nRfl.Value().FieldByName(model.KeyPK)
		if fldT.Type == reflect.TypeOf(uuid.UUID{}) && fld.IsZero() {
			fld.Set(reflect.ValueOf(uuid.New()))
		}
	})
	return rfl
}
