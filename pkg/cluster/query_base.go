package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryBase struct {
	_svc ServiceChain
	_qr  *QueryRequest
}

func (q *queryBase) baseInit(serviceChain ServiceChain, variant QueryVariant) {
	q._svc = serviceChain
	q._qr = &QueryRequest{Variant: variant, opts: map[string]interface{}{}}
}

func (q *queryBase) baseModel(m interface{}) {
	q._qr.Model = model.NewReflect(m)
}

func (q *queryBase) baseQueryRequest() *QueryRequest {
	return q._qr
}

func (q *queryBase) baseServiceChain() ServiceChain {
	return q._svc
}

func (q *queryBase) baseExec(ctx context.Context) error {
	return q.baseServiceChain().Exec(ctx, q.baseQueryRequest())
}
