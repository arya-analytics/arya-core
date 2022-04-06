package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Engine struct {
	pool   *pool.Pool[*Engine]
	driver Driver
}

func New(driver Driver) *Engine {
	e := &Engine{driver: driver, pool: pool.New[*Engine]()}
	e.pool.Factory = e
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
		&query.Create{}:   newCreate(c).exec,
		&query.Retrieve{}: newRetrieve(c).exec,
		&query.Delete{}:   newDelete(c).exec,
		&query.Migrate{}:  newMigrate(c).exec,
	}, query.SwitchWithoutPanic())
	e.pool.Release(a)
	return newErrorConvert().Exec(err)
}

func (e *Engine) NewAdapt(*Engine) (pool.Adapt[*Engine], error) {
	return newAdapter(e.driver)
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
	_, ok := p.Query().(*query.Migrate)
	if ok {
		return true
	}
	if !internal.RequiresEngine(p.Model(), e) {
		return false
	}
	fieldsOpt, ok := query.RetrieveFieldsOpt(p)
	if ok {
		return p.Model().StructTagChain().HasAnyFields(fieldsOpt.AllExcept("ID")...)
	}
	return true
}

func (e *Engine) NewCreate() *query.Create {
	return query.NewCreate().BindExec(e.Exec)
}

func (e *Engine) NewRetrieve() *query.Retrieve {
	return query.NewRetrieve().BindExec(e.Exec)
}

func (e *Engine) NewDelete() *query.Delete {
	return query.NewDelete().BindExec(e.Exec)
}

func (e *Engine) NewMigrate() *query.Migrate {
	return query.NewMigrate().BindExec(e.Exec)
}
