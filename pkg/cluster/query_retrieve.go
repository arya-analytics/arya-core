package cluster

import "context"

type QueryRetrieve struct {
	queryBase
}

func newRetrieve(svc ServiceChain) *QueryRetrieve {
	q := &QueryRetrieve{}
	q.baseInit(svc, QueryVariantRetrieve)
	return q
}

func (q *QueryRetrieve) Model(m interface{}) *QueryRetrieve {
	q.baseModel(m)
	return q
}

func (q *QueryRetrieve) WherePK(pk interface{}) *QueryRetrieve {
	NewPKQueryOpt(q.baseQueryRequest(), []interface{}{pk})
	return q
}

func (q *QueryRetrieve) WherePKs(pks interface{}) *QueryRetrieve {
	NewPKQueryOpt(q.baseQueryRequest(), pks)
	return q
}

func (q *QueryRetrieve) WhereFields(flds Fields) *QueryRetrieve {
	NewFieldsQueryOpt(q.baseQueryRequest(), flds)
	return q
}

func (q *QueryRetrieve) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
