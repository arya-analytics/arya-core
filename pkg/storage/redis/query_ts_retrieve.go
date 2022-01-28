package redis

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
	"time"
)

type tsRetrieveQuery struct {
	baseQuery
	PKs      []model.PK
	allRange bool
	fromTS   time.Time
	toTS     time.Time
}

func newTSRetrieve(client *timeseries.Client) *tsRetrieveQuery {
	c := &tsRetrieveQuery{}
	c.baseInit(client)
	return c
}

func (tsr *tsRetrieveQuery) Model(m interface{}) storage.CacheTSRetrieveQuery {
	tsr.baseModel(m)
	return tsr
}

func (tsr *tsRetrieveQuery) WherePK(pk interface{}) storage.CacheTSRetrieveQuery {
	tsr.PKs = append(tsr.PKs, model.NewPK(pk))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePKs(pks interface{}) storage.CacheTSRetrieveQuery {
	rfl := reflect.ValueOf(pks)
	for i := 0; i < rfl.Len(); i++ {
		tsr.WherePK(rfl.Index(i).Interface())
	}
	return tsr
}

func (tsr *tsRetrieveQuery) AllTimeRange() storage.CacheTSRetrieveQuery {
	tsr.allRange = true
	return tsr
}

func (tsr *tsRetrieveQuery) WhereTimeRange(fromTS time.Time, toTS time.Time) storage.CacheTSRetrieveQuery {
	tsr.fromTS = fromTS
	tsr.toTS = toTS
	return tsr
}

func (tsr *tsRetrieveQuery) SeriesExists(ctx context.Context, pk interface{}) (bool,
	error) {
	var res interface{}
	tsr.catcher.Exec(func() (err error) {
		res, err = tsr.baseClient().Exists(ctx, model.NewPK(pk).String()).Result()
		return err
	})
	return res.(int64) != 0, tsr.baseErr()
}

func (tsr *tsRetrieveQuery) Exec(ctx context.Context) error {
	wrapper := &tsModelWrapper{rfl: tsr.modelAdapter.Dest()}
	tsr.validateReq()
	for _, pk := range tsr.PKs {
		tsr.catcher.Exec(func() error {
			var err error
			var res interface{}
			if tsr.allRange {
				res, err = tsr.baseClient().TSGetAll(ctx, pk.String()).Result()
			} else if !tsr.fromTS.IsZero() {
				res, err = tsr.baseClient().TSGetRange(ctx, pk.String(), tsr.fromTS.UnixNano(),
					tsr.toTS.UnixNano()).Result()
			} else {
				res, err = tsr.baseClient().TSGet(ctx, pk.String()).Result()
			}
			if err != nil {
				return err
			}
			return wrapper.bindRes(pk.String(), res)
		})
	}
	tsr.baseAdaptToSource()
	return tsr.baseErr()
}

func (tsr *tsRetrieveQuery) validateReq() {
	tsr.catcher.Exec(func() error {
		return tsRetrieveQueryReqValidator.Exec(tsr)
	})
}

var tsRetrieveQueryReqValidator = validate.New([]validate.Func{
	validatePKProvided,
})

func validatePKProvided(v interface{}) error {
	q := v.(*tsRetrieveQuery)
	if (len(q.PKs)) == 0 {
		return storage.Error{
			Type:    storage.ErrTypeInvalidArgs,
			Message: fmt.Sprintf("no PK provided to ts retrieve query"),
		}
	}
	return nil
}
