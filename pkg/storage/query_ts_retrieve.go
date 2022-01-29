package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type tsRetrieveQuery struct {
	tsBaseQuery
	modelRfl *model.Reflect
}

func newTSRetrieve(s *Storage) *tsRetrieveQuery {
	tsr := &tsRetrieveQuery{}
	tsr.baseInit(s)
	return tsr
}

func (tsr *tsRetrieveQuery) Model(m interface{}) *tsRetrieveQuery {
	tsr.tsBaseModel(m)
	tsr.setCacheQuery(tsr.cacheQuery().Model(m))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePK(pk interface{}) *tsRetrieveQuery {
	tsr.tsBaseWherePk(pk)
	tsr.setCacheQuery(tsr.cacheQuery().WherePK(pk))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePKs(pks interface{}) *tsRetrieveQuery {
	tsr.tsBaseWherePks(pks)
	tsr.setCacheQuery(tsr.cacheQuery().WherePKs(pks))
	return tsr
}

func (tsr *tsRetrieveQuery) AllTimeRange() *tsRetrieveQuery {
	tsr.setCacheQuery(tsr.cacheQuery().AllTimeRange())
	return tsr
}

func (tsr *tsRetrieveQuery) WhereTimeRange(fromTS int64, toTS int64) *tsRetrieveQuery {
	tsr.setCacheQuery(tsr.cacheQuery().WhereTimeRange(fromTS, toTS))
	return tsr
}

func (tsr *tsRetrieveQuery) SeriesExists(ctx context.Context, pk interface{}) (bool, error) {
	return tsr.cacheQuery().SeriesExists(ctx, pk)
}

func (tsr *tsRetrieveQuery) Exec(ctx context.Context) error {
	tsr.catcher.Exec(func() error {
		err := tsr.cacheQuery().Exec(ctx)
		if err != nil {
			se := err.(Error)
			if se.Type == ErrTypeItemNotFound {
				if bErr := tsr.tsBaseCreateIndexes(ctx, tsr.pks); bErr != nil {
					return bErr
				}
				// retry the transaction after we've created the indexes
				return tsr.Exec(ctx)
			}
		}
		return err
	})
	return tsr.baseErr()
}

// |||| QUERY BINDING |||

// || CACHE ||

func (tsr *tsRetrieveQuery) cacheQuery() CacheTSRetrieveQuery {
	if tsr.baseCacheQuery() == nil {
		tsr.setCacheQuery(tsr.cacheEngine.NewTSRetrieve(tsr.baseCacheAdapter()))
	}
	return tsr.baseCacheQuery().(CacheTSRetrieveQuery)
}

func (tsr *tsRetrieveQuery) setCacheQuery(q CacheTSRetrieveQuery) {
	tsr.baseSetCacheQuery(q)
}
