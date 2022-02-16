package cluster

import "context"

type QueryDelete struct {
	queryBase
}

func newDelete(svc ServiceChain) *QueryDelete {
	q := &QueryDelete{}
	q.baseInit(svc, QueryVariantDelete)
	return q
}

func (q *QueryDelete) Model(m interface{}) *QueryDelete {
	q.baseModel(m)
	return q
}

func (q *QueryDelete) WherePK(pk interface{}) *QueryDelete {
	NewPKQueryOpt(q.baseQueryRequest(), pk)
	return q
}

func (q *QueryDelete) WherePKs(pks interface{}) *QueryDelete {
	NewPKQueryOpt(q.baseQueryRequest(), pks)
	return q
}

func (q *QueryDelete) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
