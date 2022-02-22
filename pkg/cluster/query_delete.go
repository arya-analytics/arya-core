package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type QueryDelete struct {
	queryBase
}

func newDelete(svc ServiceChain) *QueryDelete {
	q := &QueryDelete{}
	q.baseInit(svc)
	return q
}

func (q *QueryDelete) Model(m interface{}) *QueryDelete {
	q.baseModel(internal.QueryVariantDelete, m)
	return q
}

func (q *QueryDelete) WherePK(pk interface{}) *QueryDelete {
	internal.NewPKQueryOpt(q.baseQueryRequest(), []interface{}{pk})
	return q
}

func (q *QueryDelete) WherePKs(pks interface{}) *QueryDelete {
	internal.NewPKQueryOpt(q.baseQueryRequest(), pks)
	return q
}

func (q *QueryDelete) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
