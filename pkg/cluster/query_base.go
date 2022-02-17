package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryBase struct {
	_svc ServiceChain
	_qr  *internal.QueryRequest
}

func (q *queryBase) baseInit(serviceChain ServiceChain) {
	q._svc = serviceChain
}

func (q *queryBase) baseModel(variant internal.QueryVariant, m interface{}) {
	q._qr = internal.NewQueryRequest(variant, model.NewReflect(m))
}

func (q *queryBase) baseQueryRequest() *internal.QueryRequest {
	return q._qr
}

func (q *queryBase) baseServiceChain() ServiceChain {
	return q._svc
}

func (q *queryBase) baseExec(ctx context.Context) error {
	return q.baseServiceChain().Exec(ctx, q.baseQueryRequest())
}
