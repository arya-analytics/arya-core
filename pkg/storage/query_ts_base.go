package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryTSBase struct {
	pks model.PKChain
	queryBase
}

func (q *queryTSBase) tsBaseWherePk(pk interface{}) {
	q.pks = append(q.pks, model.NewPK(pk))
}

func (q *queryTSBase) tsBaseWherePks(pks interface{}) {
	q.pks = model.NewPKChain(pks)
}
