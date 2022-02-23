package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Engine struct {
	driver Driver
}

func New(driver Driver) *Engine {
	return &Engine{driver}
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

func (e *Engine) NewCreate(a storage.Adapter) storage.QueryObjectCreate {
	return newCreate(conn(a))
}

func (e *Engine) NewRetrieve(a storage.Adapter) storage.QueryObjectRetrieve {
	return newRetrieve(conn(a))
}

func (e *Engine) NewDelete(a storage.Adapter) storage.QueryObjectDelete {
	return newDelete(conn(a))
}

func (e *Engine) NewMigrate(a storage.Adapter) storage.QueryObjectMigrate {
	return newMigrate(conn(a))
}
