package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/go-redis/redis/v8"
)

type tsRetrieveQuery struct {
	baseQuery
	PKChain  model.PKChain
	allRange bool
	fromTS   int64
	toTS     int64
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
	tsr.PKChain = append(tsr.PKChain, model.NewPK(pk))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePKs(pks interface{}) storage.CacheTSRetrieveQuery {
	tsr.PKChain = model.NewPKChain(pks)
	return tsr
}

func (tsr *tsRetrieveQuery) AllTimeRange() storage.CacheTSRetrieveQuery {
	tsr.allRange = true
	return tsr
}

func (tsr *tsRetrieveQuery) WhereTimeRange(fromTS int64, toTS int64) storage.CacheTSRetrieveQuery {
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
	tsr.validateReq()
	for _, pk := range tsr.PKChain {
		tsr.catcher.Exec(func() error {
			pks := pk.String()
			var cmd *redis.Cmd
			if tsr.allRange {
				cmd = tsr.baseClient().TSGetAll(ctx, pks)
			} else if tsr.toTS != 0 {
				cmd = tsr.baseClient().TSGetRange(ctx, pks, tsr.fromTS, tsr.toTS)
			} else {
				cmd = tsr.baseClient().TSGet(ctx, pk.String())
			}
			res, err := cmd.Result()
			if err != nil {
				return err
			}
			return tsr.modelAdapter.bindRes(pk.String(), res)
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
	if (len(q.PKChain)) == 0 {
		return storage.Error{
			Type:    storage.ErrTypeInvalidArgs,
			Message: "no PK provided to ts retrieve query",
		}
	}
	return nil
}
