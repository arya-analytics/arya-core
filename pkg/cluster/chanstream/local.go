package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

type ServiceLocalStorage struct {
	store storage.Storage
}

func (s *ServiceLocalStorage) create(ctx context.Context, p *query.Pack) error {
	goExecOpt, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}
	for {
		sample, sampleOK := p.Model().ChanRecv()
		if !sampleOK {
			break
		}
		if err := tsquery.NewCreate().Model(sample).Exec(ctx); err != nil {
			goExecOpt.Errors <- err
		}
	}
	return nil
}

func (s *ServiceLocalStorage) retrieve(ctx context.Context, p *query.Pack) error {

}
