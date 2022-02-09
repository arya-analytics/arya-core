package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryBase struct {
	storage         Storage
	modelRfl        *model.Reflect
	_baseMDQuery    QueryMDBase
	_baseObjQuery   QueryObjectBase
	_baseCacheQuery QueryCacheBase
	_catcher        *errutil.Catcher
}

func (q *queryBase) baseInit(s Storage) {
	q.storage = s
	q._catcher = &errutil.Catcher{}
}

// |||| MODEL UTILITIES ||||

func (q *queryBase) baseBindModel(m interface{}) {
	q.modelRfl = model.NewReflect(m)
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (q *queryBase) baseMDEngine() EngineMD {
	return q.storage.config().EngineMD
}

func (q *queryBase) baseMDAdapter() Adapter {
	return q.storage.adapter(q.baseMDEngine())
}

func (q *queryBase) baseMDQuery() QueryMDBase {
	return q._baseMDQuery
}

func (q *queryBase) baseSetMDQuery(qmd QueryMDBase) {
	q._baseMDQuery = qmd
}

// || OBJECT ||

func (q *queryBase) baseObjEngine() EngineObject {
	return q.storage.config().EngineObject
}

func (q *queryBase) baseObjAdapter() Adapter {
	return q.storage.adapter(q.baseObjEngine())
}

func (q *queryBase) baseObjQuery() QueryObjectBase {
	return q._baseObjQuery
}

func (q *queryBase) baseSetObjQuery(qob QueryObjectBase) {
	q._baseObjQuery = qob
}

// || CACHE ||

func (q *queryBase) baseCacheEngine() EngineCache {
	return q.storage.config().EngineCache
}

func (q *queryBase) baseCacheAdapter() Adapter {
	return q.storage.adapter(q.baseCacheEngine())
}

func (q *queryBase) baseCacheQuery() QueryCacheBase {
	return q._baseCacheQuery
}

func (q *queryBase) baseSetCacheQuery(qca QueryCacheBase) {
	q._baseCacheQuery = qca
}

// |||| EXCEPTION HANDLING  ||||

func (q *queryBase) baseExec(actionFunc errutil.ActionFunc) {
	q._catcher.Exec(actionFunc)
}

func (q *queryBase) baseErr() error {
	return q._catcher.Error()
}
