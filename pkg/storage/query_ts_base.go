package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type tsBaseQuery struct {
	pks model.PKChain
	baseQuery
}

func (tsb *tsBaseQuery) tsBaseWherePk(pk interface{}) {
	tsb.pks = append(tsb.pks, model.NewPK(pk))
}

func (tsb *tsBaseQuery) tsBaseWherePks(pks interface{}) {
	tsb.pks = model.NewPKChain(pks)
}
