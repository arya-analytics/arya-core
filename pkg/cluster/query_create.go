package cluster

import "context"

type QueryCreate struct {
	queryBase
}

func newCreate(svc ServiceChain) *QueryCreate {
	q := &QueryCreate{}
	q.baseInit(svc, QueryVariantCreate)
	return q
}

func (q *QueryCreate) Model(m interface{}) *QueryCreate {
	q.baseModel(m)
	return q
}

func (q *QueryCreate) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
