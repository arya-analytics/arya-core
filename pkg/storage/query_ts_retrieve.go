package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type tsRetrieveQuery struct {
	baseQuery
	pks      []interface{}
	modelRfl *model.Reflect
}

func newTSRetrieve(s *Storage) *tsRetrieveQuery {
	tsr := &tsRetrieveQuery{}
	tsr.baseInit(s)
	return tsr
}

func (tsr *tsRetrieveQuery) Model(m interface{}) *tsRetrieveQuery {
	tsr.modelRfl = model.NewReflect(m)
	tsr.setCacheQuery(tsr.cacheQuery().Model(m))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePK(pk interface{}) *tsRetrieveQuery {
	tsr.pks = append(tsr.pks, pk)
	tsr.setCacheQuery(tsr.cacheQuery().WherePK(pk))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePKs(pks interface{}) *tsRetrieveQuery {
	rv := reflect.ValueOf(pks)
	for i := 0; i < rv.Len(); i++ {
		tsr.pks = append(tsr.pks, rv.Index(i))
	}
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
		se, ok := err.(Error)
		if !ok {
			panic(err)
		}
		if se.Type == ErrTypeItemNotFound {
			log.Warn("Series not found in cache. Attempting to fix.")
			r := newRetrieve(tsr.storage)
			tag, ok := tsr.modelRfl.Tags().Retrieve("storage", "role", "index")
			if !ok {
				return err
			}
			fld, ok := tsr.modelRfl.Type().FieldByName(tag.FldName)
			if !ok {
				return err
			}
			sm := catalog().NewFromType(fld.Type.Elem(), true)
			if sErr := r.Model(sm).WherePKs(tsr.pks).Exec(ctx); err != nil {
				return sErr
			}
			if tscErr := newTSCreate(tsr.storage).Model(sm).Series().Exec(
				ctx); tscErr != nil {
				return tscErr
			}
			// retry the transaction after we've created the indexes
			return tsr.Exec(ctx)
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
