package route

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

func BatchModel[T comparable](m interface{}, fld string) map[T]*model.Reflect {
	rfl := model.NewReflect(m)
	b := map[T]*model.Reflect{}
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		rawFldV := nRfl.StructFieldByName(fld)
		fldV, ok := rawFldV.Interface().(T)
		if !ok {
			panic(fmt.Sprintf("batch model received unknown type for field. received type %s", rawFldV.Type()))
		}
		v, ok := b[fldV]
		if !ok {
			v = nRfl.NewChain()
			b[fldV] = v
		}
		v.ChainAppend(nRfl)
	})
	return b
}
