package storage

import "context"

type QueryTSCreate struct {
	queryTSBase
}

func newTSCreate(s Storage) *QueryTSCreate {
	q := &QueryTSCreate{}
	q.baseInit(s, s.config().Hooks.BeforeTSCreate)
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

func (q *QueryTSCreate) cacheQuery() QueryCacheTSCreate {
	if q.baseCacheQuery() == nil {
		q.setCacheQuery(q.baseCacheEngine().NewTSCreate(q.baseCacheAdapter()))
	}
	return q.baseCacheQuery().(QueryCacheTSCreate)
}

func (q *QueryTSCreate) setCacheQuery(qca QueryCacheTSCreate) {
	q.baseSetCacheQuery(qca)
}
