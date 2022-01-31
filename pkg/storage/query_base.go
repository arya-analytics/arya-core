package storage

import "github.com/arya-analytics/aryacore/pkg/util/errutil"

type baseQuery struct {
	storage         *Storage
	mdEngine        MDEngine
	_baseMDQuery    MDBaseQuery
	objEngine       ObjectEngine
	_baseObjQuery   ObjectBaseQuery
	cacheEngine     CacheEngine
	_baseCacheQuery CacheBaseQuery
	catcher         *errutil.Catcher
}

func (b *baseQuery) baseInit(s *Storage) {
	b.storage = s
	b.mdEngine = b.storage.cfg.mdEngine()
	b.objEngine = b.storage.cfg.objEngine()
	b.cacheEngine = b.storage.cfg.cacheEngine()
	b.catcher = &errutil.Catcher{}
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (b *baseQuery) baseMDAdapter() Adapter {
	return b.storage.adapter(EngineRoleMD)
}

func (b *baseQuery) baseMDQuery() MDBaseQuery {
	return b._baseMDQuery
}

func (b *baseQuery) baseSetMDQuery(q MDBaseQuery) {
	b._baseMDQuery = q
}

// || OBJECT ||

func (b *baseQuery) baseObjAdapter() Adapter {
	return b.storage.adapter(EngineRoleObject)
}

func (b *baseQuery) baseObjQuery() ObjectBaseQuery {
	return b._baseObjQuery
}

func (b *baseQuery) baseSetObjQuery(q ObjectBaseQuery) {
	b._baseObjQuery = q
}

// || CACHE ||

func (b *baseQuery) baseCacheAdapter() Adapter {
	return b.storage.adapter(EngineRoleCache)
}

func (b *baseQuery) baseCacheQuery() CacheBaseQuery {
	return b._baseCacheQuery
}

func (b *baseQuery) baseSetCacheQuery(q CacheBaseQuery) {
	b._baseCacheQuery = q
}

// |||| EXCEPTION HANDLING  ||||

func (b *baseQuery) baseErr() error {
	return b.catcher.Error()
}
