package storage

import "context"

type TSCreateQuery struct {
	tsBaseQuery
}

func newTSCreate(s *Storage) *TSCreateQuery {
	tsc := &TSCreateQuery{}
	tsc.baseInit(s)
	return tsc
}

func (tsc *TSCreateQuery) Model(m interface{}) *TSCreateQuery {
	tsc.baseBindModel(m)
	tsc.setCacheQuery(tsc.cacheQuery().Model(m))
	return tsc
}

func (tsc *TSCreateQuery) Series() *TSCreateQuery {
	tsc.setCacheQuery(tsc.cacheQuery().Series())
	return tsc
}

func (tsc *TSCreateQuery) Sample() *TSCreateQuery {
	tsc.setCacheQuery(tsc.cacheQuery().Sample())
	return tsc
}

func (tsc *TSCreateQuery) Exec(ctx context.Context) error {
	tsc.baseExec(func() error { return tsc.cacheQuery().Exec(ctx) })
	return tsc.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (tsc *TSCreateQuery) cacheQuery() CacheTSCreateQuery {
	if tsc.baseCacheQuery() == nil {
		tsc.setCacheQuery(tsc.baseCacheEngine().NewTSCreate(tsc.baseCacheAdapter()))
	}
	return tsc.baseCacheQuery().(CacheTSCreateQuery)
}

func (tsc *TSCreateQuery) setCacheQuery(q CacheTSCreateQuery) {
	tsc.baseSetCacheQuery(q)
}
