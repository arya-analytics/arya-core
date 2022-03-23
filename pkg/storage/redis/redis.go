package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

// |||| ENGINE ||||

type Engine struct {
	pool   *storage.Pool
	driver Driver
	tsquery.AssembleTS
}

func New(driver Driver, pool *storage.Pool) *Engine {
	e := &Engine{driver: driver, pool: pool}
	e.AssembleTS = tsquery.NewAssemble(e.Exec)
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
		&tsquery.Create{}:   newTSCreate(c).exec,
		&tsquery.Retrieve{}: newTSRetrieve(c).exec,
	})
	e.pool.Release(a)
	return newErrorConvert().Exec(err)
}

func (e *Engine) NewAdapter() (internal.Adapter, error) {
	return newAdapter(e.driver)
}

func (e *Engine) IsAdapter(a internal.Adapter) bool {
	_, ok := a.(*adapter)
	return ok
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
	switch p.Query().(type) {
	case *tsquery.Retrieve:
		return true
	case *tsquery.Create:
		return true
	default:
		return false
	}
}
