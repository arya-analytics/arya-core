package storage

import "context"

type tsCreateQuery struct {
	tsBaseQuery
}

func newTSCreate(s *Storage) *tsCreateQuery {
	tsc := &tsCreateQuery{}
	tsc.baseInit(s)
	return tsc
}

func (tsc *tsCreateQuery) Model(m interface{}) *tsCreateQuery {
	tsc.setCacheQuery(tsc.cacheQuery().Model(m))
	return tsc
}

func (tsc *tsCreateQuery) Series() *tsCreateQuery {
	tsc.setCacheQuery(tsc.cacheQuery().Series())
	return tsc
}

func (tsc *tsCreateQuery) Sample() *tsCreateQuery {
	tsc.setCacheQuery(tsc.cacheQuery().Sample())
	return tsc
}

func (tsc *tsCreateQuery) Exec(ctx context.Context) error {
	tsc.catcher.Exec(func() error {
		return tsc.cacheQuery().Exec(ctx)
	})
	return tsc.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (tsc *tsCreateQuery) cacheQuery() CacheTSCreateQuery {
	if tsc.baseCacheQuery() == nil {
		tsc.setCacheQuery(tsc.storage.cfg.cacheEngine().NewTSCreate(tsc.baseCacheAdapter()))
	}
	return tsc.baseCacheQuery().(CacheTSCreateQuery)
}

func (tsc *tsCreateQuery) setCacheQuery(q CacheTSCreateQuery) {
	tsc.baseSetCacheQuery(q)
}
