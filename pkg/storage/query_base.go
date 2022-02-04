package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type baseQuery struct {
	storage         *Storage
	modelRfl        *model.Reflect
	_catcher        *errutil.Catcher
	_baseMDQuery    MDBaseQuery
	_baseObjQuery   ObjectBaseQuery
	_baseCacheQuery CacheBaseQuery
}

func (b *baseQuery) baseInit(s *Storage) {
	b.storage = s
	b._catcher = &errutil.Catcher{}
}

// |||| MODEL UTILITIES ||||

func (b *baseQuery) baseBindModel(m interface{}) {
	b.modelRfl = model.NewReflect(m)
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (b *baseQuery) baseMDEngine() MDEngine {
	return b.storage.cfg.MDEngine

}

func (b *baseQuery) baseMDAdapter() Adapter {
	return b.storage.adapter(b.baseMDEngine())
}

func (b *baseQuery) baseMDQuery() MDBaseQuery {
	return b._baseMDQuery
}

func (b *baseQuery) baseSetMDQuery(q MDBaseQuery) {
	b._baseMDQuery = q
}

// || OBJECT ||

func (b *baseQuery) baseObjEngine() ObjectEngine {
	return b.storage.cfg.ObjectEngine
}

func (b *baseQuery) baseObjAdapter() Adapter {
	return b.storage.adapter(b.baseObjEngine())
}

func (b *baseQuery) baseObjQuery() ObjectBaseQuery {
	return b._baseObjQuery
}

func (b *baseQuery) baseSetObjQuery(q ObjectBaseQuery) {
	b._baseObjQuery = q
}

// || CACHE ||

func (b *baseQuery) baseCacheEngine() CacheEngine {
	return b.storage.cfg.CacheEngine
}

func (b *baseQuery) baseCacheAdapter() Adapter {
	return b.storage.adapter(b.baseCacheEngine())
}

func (b *baseQuery) baseCacheQuery() CacheBaseQuery {
	return b._baseCacheQuery
}

func (b *baseQuery) baseSetCacheQuery(q CacheBaseQuery) {
	b._baseCacheQuery = q
}

// |||| EXCEPTION HANDLING  ||||

func (b *baseQuery) baseExec(actionFunc errutil.ActionFunc) {
	b._catcher.Exec(actionFunc)
}

func (b *baseQuery) baseErr() error {
	return b._catcher.Error()
}
