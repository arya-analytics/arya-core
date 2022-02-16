package batch

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Model struct {
	*model.Reflect
}

func NewModel(rfl *model.Reflect) *Model {
	return &Model{Reflect: rfl}

}

func (m *Model) Exec(key string) map[interface{}]*model.Reflect {
	b := map[interface{}]*model.Reflect{}
	m.ForEach(func(rfl *model.Reflect, i int) {
		fldV := rfl.StructFieldByName(key).Interface()
		v, ok := b[fldV]
		if !ok {
			v = rfl.NewChain()
			b[fldV] = v
		}
		v.ChainAppend(rfl)
	})
	return b
}
