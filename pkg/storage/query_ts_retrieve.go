package storage

import (
	"context"
)

type QueryTSRetrieve struct {
	queryTSBase
}

func newTSRetrieve(s *Storage) *QueryTSRetrieve {
	q := &QueryTSRetrieve{}
	q.baseInit(s)
	return q
}

func (q *QueryTSRetrieve) Model(m interface{}) *QueryTSRetrieve {
	q.baseBindModel(m)
	q.setCacheQuery(q.cacheQuery().Model(m))
	return q
}

func (q *QueryTSRetrieve) WherePK(pk interface{}) *QueryTSRetrieve {
	q.tsBaseWherePk(pk)
	q.setCacheQuery(q.cacheQuery().WherePK(pk))
	return q
}

func (q *QueryTSRetrieve) WherePKs(pks interface{}) *QueryTSRetrieve {
	q.tsBaseWherePks(pks)
	q.setCacheQuery(q.cacheQuery().WherePKs(pks))
	return q
}

func (q *QueryTSRetrieve) AllTimeRange() *QueryTSRetrieve {
	q.setCacheQuery(q.cacheQuery().AllTimeRange())
	return q
}

func (q *QueryTSRetrieve) WhereTimeRange(fromTS int64, toTS int64) *QueryTSRetrieve {
	q.setCacheQuery(q.cacheQuery().WhereTimeRange(fromTS, toTS))
	return q
}

func (q *QueryTSRetrieve) SeriesExists(ctx context.Context,
	pk interface{}) (exists bool, err error) {
	q.baseExec(func() error {
		exists, err = q.cacheQuery().SeriesExists(ctx, pk)
		return err
	})
	return exists, q.baseErr()
}

func (q *QueryTSRetrieve) Exec(ctx context.Context) error {
	q.baseExec(func() error { return q.cacheQuery().Exec(ctx) })
	return q.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (q *QueryTSRetrieve) cacheQuery() QueryCacheTSRetrieve {
	if q.baseCacheQuery() == nil {
		q.setCacheQuery(q.baseCacheEngine().NewTSRetrieve(q.baseCacheAdapter()))
	}
	return q.baseCacheQuery().(QueryCacheTSRetrieve)
}

func (q *QueryTSRetrieve) setCacheQuery(qca QueryCacheTSRetrieve) {
	q.baseSetCacheQuery(qca)
}
