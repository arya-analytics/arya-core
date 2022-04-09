package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"time"
)

type Persist struct {
	*mock.DataSourceMem
}

func (ps *Persist) Exec(ctx context.Context, p *query.Pack) error {
	switch p.Query().(type) {
	case *streamq.TSCreate:
		if p.Model().IsChan() {
			return ps.tsChanCreate(ctx, p)
		} else {
			return ps.tsCreate(ctx, p)
		}
	case *streamq.TSRetrieve:
		if p.Model().IsChan() {
			return ps.tsChanRetrieve(ctx, p)
		} else {
			return ps.tsRetrieve(ctx, p)
		}
	default:
		return ps.DataSourceMem.Exec(ctx, p)
	}
	return nil
}

func (ps *Persist) tsCreate(ctx context.Context, p *query.Pack) error {
	return ps.DataSourceMem.NewCreate().Model(p.Model()).Exec(ctx)
}

func (ps *Persist) tsChanCreate(ctx context.Context, p *query.Pack) error {
	goe, _ := streamq.RetrieveStreamOpt(p)
	go func() {
		for {
			rfl, ok := p.Model().ChanRecv()
			if !ok {
				break
			}
			if err := ps.NewCreate().Model(rfl).Exec(ctx); err != nil {
				goe.Errors <- err
			}
		}
	}()
	return nil
}

const clockSpeed = 10 * time.Millisecond

func (ps *Persist) tsRetrieve(ctx context.Context, p *query.Pack) error {
	pkc, ok := query.RetrievePKOpt(p)
	if !ok {
		panic("query must have a pk specified")
	}
	for _, pk := range pkc {
		s := &models.ChannelSample{
			ChannelConfigID: pk.Raw().(uuid.UUID),
			Timestamp:       telem.NewTimeStamp(time.Now()),
			Value:           rand.Float64(),
		}
		if p.Model().IsStruct() {
			p.Model().StructValue().Set(reflect.ValueOf(s))
		} else {
			p.Model().ChainAppend(model.NewReflect(s))
		}
	}
	return nil
}

func (ps *Persist) tsChanRetrieve(ctx context.Context, p *query.Pack) error {
	t := time.NewTicker(clockSpeed)
	pkOpt, ok := query.RetrievePKOpt(p)
	if !ok {
		panic("query must have a pk specified")
	}
	go func() {
		for range t.C {
			for _, pk := range pkOpt {
				p.Model().ChanSend(model.NewReflect(&models.ChannelSample{
					ChannelConfigID: pk.Raw().(uuid.UUID),
					Timestamp:       telem.NewTimeStamp(time.Now()),
					Value:           rand.Float64(),
				}))
			}
		}
	}()
	return nil
}
