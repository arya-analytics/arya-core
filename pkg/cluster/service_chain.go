package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type ServiceChain []Service

func (sc ServiceChain) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	for _, s := range sc {
		if s.CanHandle(qr) {
			return s.Exec(ctx, qr)
		}
	}
	panic("no service could handle the request!")
}
