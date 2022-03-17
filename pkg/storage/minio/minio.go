package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/minio/minio-go/v7"
)

type Engine struct {
	pool   *storage.Pool
	driver Driver
}

func New(driver Driver, pool *storage.Pool) *Engine {
	return &Engine{driver: driver, pool: pool}
}

func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	if !e.shouldHandle(p) {
		return nil
	}

	c, err := e.client()
	if err != nil {
		return newErrorConvert().Exec(err)
	}

	return query.Switch(ctx, p, query.Ops{
		Create:   newCreate(c).exec,
		Retrieve: newRetrieve(c).exec,
		Delete:   newDelete(c).exec,
		Migrate:  newMigrate(c).exec,
	})
}

func (e *Engine) client() (*minio.Client, error) {
	a, err := e.pool.Retrieve(e)
	if err != nil {
		return nil, err
	}
	return conn(a), nil

}

func (e *Engine) NewAdapter() (internal.Adapter, error) {
	return newAdapter(e.driver)
}

func (e *Engine) IsAdapter(a internal.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
	_, ok := p.Query().(*query.Migrate)
	if ok {
		return true
	}
	if !catalog().Contains(p.Model().Pointer()) {
		return false
	}
	fldsOpt, ok := query.RetrieveFieldsOpt(p)
	if ok {
		rfl := model.NewReflect(catalog().New(p.Model().Pointer()))
		return rfl.StructTagChain().HasAnyFields(fldsOpt.AllExcept("ID")...)
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
