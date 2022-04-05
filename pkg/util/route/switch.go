package route

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

func ModelSwitchIter[T comparable](m interface{}, fld string, action func(fld T, rfl *model.Reflect)) {
	switchMap := BatchModel[T](m, fld)
	for k, v := range switchMap {
		action(k, v)
	}
}

func ModelSwitchBoolean(m interface{}, boolFld string, trueAction, falseAction func(rfl *model.Reflect)) {
	switchMap := BatchModel[bool](m, boolFld)
	if trueRfl, ok := switchMap[true]; ok {
		trueAction(trueRfl)
	}
	if falseRfl, ok := switchMap[false]; ok {
		falseAction(falseRfl)
	}
}
