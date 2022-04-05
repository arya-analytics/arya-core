package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type ServiceChain []Service

func (sc ServiceChain) Exec(ctx context.Context, p *query.Pack) error {
	for _, s := range sc {
		if s.CanHandle(p) {
			return s.Exec(ctx, p)
		}
	}
	panic("cluster - no service could handle the request!")
}
