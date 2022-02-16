package cluster

import "context"

type ServiceChain []Service

func (sc ServiceChain) Exec(ctx context.Context, qr *QueryRequest) error {
	for _, s := range sc {
		if s.CanHandle(qr) {
			return s.Exec(ctx, qr)
		}
	}
	panic("no service could handle the request!")
}
