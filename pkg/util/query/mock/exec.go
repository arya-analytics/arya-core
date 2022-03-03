package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Exec struct {
	Pack *query.Pack
}

func (e *Exec) Exec(ctx context.Context, p *query.Pack) error {
	e.Pack = p
	return nil
}
