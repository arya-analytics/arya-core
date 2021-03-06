package roach

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
)

// || BINDING ||

func beforeInsertSetUUID(rfl *model.Reflect) *model.Reflect {
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		pkf := nRfl.PKField()
		if pkf.IsZero() {
			switch pkf.Interface().(type) {
			case uuid.UUID:
				pkf.Set(reflect.ValueOf(uuid.New()))
			}
		}
	})
	return rfl
}
