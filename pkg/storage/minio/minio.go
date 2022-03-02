package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
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
	return query.Switch(ctx, p, query.Ops{
		Create:   newCreate(e.client()).exec,
		Retrieve: newRetrieve(e.client()).exec,
		Delete:   newDelete(e.client()).exec,
	})
}

func (e *Engine) client() *minio.Client {
	return conn(e.pool.Retrieve(e))

}

func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.driver)
}

func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
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

func (e *Engine) NewMigrate() storage.QueryObjectMigrate {
	return newMigrate(e.client())
}
