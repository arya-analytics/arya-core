package storage

import (
	"context"
)

type TSRetrieveQuery struct {
	tsBaseQuery
}

func newTSRetrieve(s *Storage) *TSRetrieveQuery {
	tsr := &TSRetrieveQuery{}
	tsr.baseInit(s)
	return tsr
}

func (tsr *TSRetrieveQuery) Model(m interface{}) *TSRetrieveQuery {
	tsr.tsBaseModel(m)
	tsr.setCacheQuery(tsr.cacheQuery().Model(m))
	return tsr
}

func (tsr *TSRetrieveQuery) WherePK(pk interface{}) *TSRetrieveQuery {
	tsr.tsBaseWherePk(pk)
	tsr.setCacheQuery(tsr.cacheQuery().WherePK(pk))
	return tsr
}

func (tsr *TSRetrieveQuery) WherePKs(pks interface{}) *TSRetrieveQuery {
	tsr.tsBaseWherePks(pks)
	tsr.setCacheQuery(tsr.cacheQuery().WherePKs(pks))
	return tsr
}

func (tsr *TSRetrieveQuery) AllTimeRange() *TSRetrieveQuery {
	tsr.setCacheQuery(tsr.cacheQuery().AllTimeRange())
	return tsr
}

func (tsr *TSRetrieveQuery) WhereTimeRange(fromTS int64, toTS int64) *TSRetrieveQuery {
	tsr.setCacheQuery(tsr.cacheQuery().WhereTimeRange(fromTS, toTS))
	return tsr
}

func (tsr *TSRetrieveQuery) SeriesExists(ctx context.Context, pk interface{}) (bool, error) {
	return tsr.cacheQuery().SeriesExists(ctx, pk)
}

func (tsr *TSRetrieveQuery) Exec(ctx context.Context) error {
	tsr.catcher.Exec(func() error {
		return tsr.cacheQuery().Exec(ctx)
	})
	return tsr.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (tsr *TSRetrieveQuery) cacheQuery() CacheTSRetrieveQuery {
	if tsr.baseCacheQuery() == nil {
		tsr.setCacheQuery(tsr.baseCacheEngine().NewTSRetrieve(tsr.baseCacheAdapter()))
	}
	return tsr.baseCacheQuery().(CacheTSRetrieveQuery)
}

func (tsr *TSRetrieveQuery) setCacheQuery(q CacheTSRetrieveQuery) {
	tsr.baseSetCacheQuery(q)
}
