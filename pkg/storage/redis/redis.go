package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
)

// |||| ENGINE ||||

type Engine struct {
	pool   *pool.Pool[internal.Engine]
	driver Driver
	streamq.AssembleTS
}

func New(driver Driver, pool *pool.Pool[internal.Engine]) *Engine {
	e := &Engine{driver: driver, pool: pool}
	e.AssembleTS = streamq.NewAssemble(e.Exec)
	return e
}

func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	if !e.shouldHandle(p) {
		return nil
	}
	a, err := e.pool.Acquire(e)
	if err != nil {
		return newErrorConvert().Exec(err)
	}
	c := client(a)
	err = query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   newTSCreate(c).exec,
		&streamq.TSRetrieve{}: newTSRetrieve(c).exec,
	})
	e.pool.Release(a)
	return newErrorConvert().Exec(err)
}

func (e *Engine) NewAdapt() (pool.Adapt[internal.Engine], error) {
	return newAdapter(e.driver)
}

func (e *Engine) Match(ce internal.Engine) bool {
	_, ok := ce.(*Engine)
	return ok
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
	switch p.Query().(type) {
	case *streamq.TSRetrieve:
		return true
	case *streamq.TSCreate:
		return true
	default:
		return false
	}
}
