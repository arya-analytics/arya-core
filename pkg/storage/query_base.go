package storage

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryBase struct {
	storage         *storage
	modelRfl        *model.Reflect
	_query          internal.Query
	_baseCacheQuery internal.QueryCacheBase
	_catcher        *errutil.CatchSimple
}

func (q *queryBase) baseInit(s *storage, query internal.Query) {
	q.storage = s
	q._query = query
	q._catcher = errutil.NewCatchSimple()
}

// |||| MODEL UTILITIES ||||

func (q *queryBase) baseBindModel(m interface{}) {
	q.modelRfl = model.NewReflect(m)
}

// || CACHE ||

func (q *queryBase) baseCacheEngine() internal.EngineCache {
	return q.storage.cfg.EngineCache
}

func (q *queryBase) baseCacheQuery() internal.QueryCacheBase {
	return q._baseCacheQuery
}

func (q *queryBase) baseSetCacheQuery(qca internal.QueryCacheBase) {
	q._baseCacheQuery = qca
}

// |||| EXCEPTION HANDLING  ||||

func (q *queryBase) baseExec(actionFunc errutil.CatchAction) {
	q._catcher.Exec(actionFunc)
}

func (q *queryBase) baseErr() error {
	return q._catcher.Error()
}
