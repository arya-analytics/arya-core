package route

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

func ModelSwitchIter[T comparable](rfl *model.Reflect, fld string, action func(fld T, rfl *model.Reflect)) {
	switchMap := BatchModel[T](rfl, fld)
	for k, v := range switchMap {
		action(k, v)
	}
}

func ModelSwitchBoolean(rfl *model.Reflect, boolFld string, trueAction, falseAction func(fld bool, rfl *model.Reflect)) {
	switchMap := BatchModel[bool](rfl, boolFld)
	if trueRfl, ok := switchMap[true]; ok {
		trueAction(true, trueRfl)
	}
	if falseRfl, ok := switchMap[false]; ok {
		falseAction(false, falseRfl)
	}
}
