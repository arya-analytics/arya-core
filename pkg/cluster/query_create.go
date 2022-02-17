package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type QueryCreate struct {
	queryBase
}

func newCreate(svc ServiceChain) *QueryCreate {
	q := &QueryCreate{}
	q.baseInit(svc)
	return q
}

func (q *QueryCreate) Model(m interface{}) *QueryCreate {
	q.baseModel(internal.QueryVariantCreate, m)
	return q
}

func (q *QueryCreate) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
