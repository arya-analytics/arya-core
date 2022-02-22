package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type QueryUpdate struct {
	queryBase
}

func newUpdate(svc ServiceChain) *QueryUpdate {
	q := &QueryUpdate{}
	q.baseInit(svc)
	return q
}

func (q *QueryUpdate) Model(m interface{}) *QueryUpdate {
	q.baseModel(internal.QueryVariantUpdate, m)
	return q
}

func (q *QueryUpdate) WherePK(pk interface{}) *QueryUpdate {
	internal.NewPKQueryOpt(q.baseQueryRequest(), []interface{}{pk})
	return q
}

func (q *QueryUpdate) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
