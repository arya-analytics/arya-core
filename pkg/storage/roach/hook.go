package roach

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
)

// || BINDING ||

func beforeInsertSetUUID(rfl *model.Reflect) *model.Reflect {
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		fldT, ok := nRfl.Type().FieldByName(pkFieldName)
		fld := nRfl.Value().FieldByName(pkFieldName)
		if !ok {
			panic(fmt.Sprintf("Detected a model with a pk field not named %s", pkFieldName))
		}
		if fldT.Type == reflect.TypeOf(uuid.UUID{}) && fld.IsZero() {
			newPK := uuid.New()
			fld.Set(reflect.ValueOf(newPK))
		}
	})
	return rfl
}
