package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

type Persist struct {
	*mock.DataSourceMem
}

func (ps *Persist) Exec(ctx context.Context, p *query.Pack) error {
	switch p.Query().(type) {
	case *tsquery.Create:
		return ps.create(ctx, p)
	case *tsquery.Retrieve:
		panic("operation not supported")
	default:
		return ps.DataSourceMem.Exec(ctx, p)
	}
	return nil
}

func (ps *Persist) create(ctx context.Context, p *query.Pack) error {
	for {
		rfl, ok := p.Model().ChanRecv()
		if !ok {
			break
		}
		ps.NewCreate().Model(rfl).Exec(ctx)
	}
	return nil
}
