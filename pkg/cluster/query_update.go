package cluster

import "context"

type QueryUpdate struct {
	queryBase
}

func newUpdate(svc ServiceChain) *QueryUpdate {
	q := &QueryUpdate{}
	q.baseInit(svc, QueryVariantUpdate)
	return q
}

func (q *QueryUpdate) Model(m interface{}) *QueryUpdate {
	q.baseModel(m)
	return q
}

func (q *QueryUpdate) WherePK(pk interface{}) *QueryUpdate {
	NewPKQueryOpt(q.baseQueryRequest(), pk)
	return q
}

func (q *QueryUpdate) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
