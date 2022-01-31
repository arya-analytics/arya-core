package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type tsBaseQuery struct {
	pks      model.PKChain
	modelRfl *model.Reflect
	baseQuery
}

func (tsb *tsBaseQuery) tsBaseModel(m interface{}) {
	tsb.modelRfl = model.NewReflect(m)
}

func (tsb *tsBaseQuery) tsBaseWherePk(pk interface{}) {
	tsb.pks = append(tsb.pks, model.NewPK(pk))
}

func (tsb *tsBaseQuery) tsBaseWherePks(pks interface{}) {
	tsb.pks = model.NewPKChain(pks)
}
