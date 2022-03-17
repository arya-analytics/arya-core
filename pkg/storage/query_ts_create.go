package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
)

type QueryTSCreate struct {
	queryTSBase
}

func newTSCreate(s *storage) *QueryTSCreate {
	q := &QueryTSCreate{}
	q.baseInit(s, q)
	return q
}

func (q *QueryTSCreate) Model(m interface{}) *QueryTSCreate {
	q.baseBindModel(m)
	q.setCacheQuery(q.cacheQuery().Model(m))
	return q
}

func (q *QueryTSCreate) Series() *QueryTSCreate {
	q.setCacheQuery(q.cacheQuery().Series())
	return q
}

func (q *QueryTSCreate) Sample() *QueryTSCreate {
	q.setCacheQuery(q.cacheQuery().Sample())
	return q
}

func (q *QueryTSCreate) Exec(ctx context.Context) error {
	q.baseExec(func() error { return q.cacheQuery().Exec(ctx) })
	return q.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (q *QueryTSCreate) cacheQuery() internal.QueryCacheTSCreate {
	if q.baseCacheQuery() == nil {
		q.setCacheQuery(q.baseCacheEngine().NewTSCreate())
	}
	return q.baseCacheQuery().(internal.QueryCacheTSCreate)
}

func (q *QueryTSCreate) setCacheQuery(qca internal.QueryCacheTSCreate) {
	q.baseSetCacheQuery(qca)
}
