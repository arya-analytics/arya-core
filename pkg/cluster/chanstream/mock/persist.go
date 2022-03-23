package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"math/rand"
	"time"
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
		if err := ps.NewCreate().Model(rfl).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

const clockSpeed = 10 * time.Millisecond

func (ps *Persist) retrieve(ctx context.Context, p *query.Pack) error {
	t := time.NewTimer(clockSpeed)
	pkOpt, ok := query.PKOpt(p)
	if !ok {
		panic("query must have a pk specified")
	}
	for range t.C {
		for _, pk := range pkOpt {
			p.Model().ChanSend(model.NewReflect(&models.ChannelSample{
				ChannelConfigID: pk.Raw().(uuid.UUID),
				Timestamp:       telem.NewTimeStamp(time.Now()),
				Value:           rand.Float64(),
			}))
		}
	}
	return nil
}
