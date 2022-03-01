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

func (e *Engine) ShouldHandle(m interface{}, flds ...string) bool {
	if !catalog().Contains(m) {
		return false
	}
	if len(flds) == 0 {
		return true
	}
	return model.NewReflect(catalog().New(m)).StructTagChain().HasAnyFields(flds...)
}

func (e *Engine) NewCreate() *query.Create {
	return query.NewCreate().BindExec(newCreate(e.client()).exec)
}

func (e *Engine) NewRetrieve() *query.Retrieve {
	return query.NewRetrieve().BindExec(newRetrieve(e.client()).exec)
}

func (e *Engine) NewDelete() *query.Delete {
	return query.NewDelete().BindExec(newDelete(e.client()).exec)
}

func (e *Engine) NewMigrate() storage.QueryObjectMigrate {
	return newMigrate(e.client())
}
