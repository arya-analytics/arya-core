package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type QueryRetrieve struct {
	queryBase
}

func newRetrieve(svc ServiceChain) *QueryRetrieve {
	q := &QueryRetrieve{}
	q.baseInit(svc)
	return q
}

func (q *QueryRetrieve) Model(m interface{}) *QueryRetrieve {
	q.baseModel(internal.QueryVariantRetrieve, m)
	return q
}

func (q *QueryRetrieve) WherePK(pk interface{}) *QueryRetrieve {
	internal.NewPKQueryOpt(q.baseQueryRequest(), pk)
	return q
}

func (q *QueryRetrieve) WherePKs(pks interface{}) *QueryRetrieve {
	internal.NewPKsQueryOpt(q.baseQueryRequest(), pks)
	return q
}

func (q *QueryRetrieve) Fields(flds ...string) *QueryRetrieve {
	internal.NewFieldsQueryOpt(q.baseQueryRequest(), flds...)
	return q
}

func (q *QueryRetrieve) Relation(rel string, flds ...string) *QueryRetrieve {
	internal.NewRelationQueryOpt(q.baseQueryRequest(), rel, flds...)
	return q
}

func (q *QueryRetrieve) WhereFields(flds model.WhereFields) *QueryRetrieve {
	internal.NewWhereFieldsQueryOpt(q.baseQueryRequest(), flds)
	return q
}

func (q *QueryRetrieve) Exec(ctx context.Context) error {
	return q.baseExec(ctx)
}
