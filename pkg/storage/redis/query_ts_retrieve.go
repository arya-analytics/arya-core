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
	queryBase
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

func (tsr *tsRetrieveQuery) Model(m interface{}) storage.QueryCacheTSRetrieve {
	tsr.baseModel(m)
	return tsr
}

func (tsr *tsRetrieveQuery) WherePK(pk interface{}) storage.QueryCacheTSRetrieve {
	tsr.PKChain = append(tsr.PKChain, model.NewPK(pk))
	return tsr
}

func (tsr *tsRetrieveQuery) WherePKs(pks interface{}) storage.QueryCacheTSRetrieve {
	tsr.PKChain = model.NewPKChain(pks)
	return tsr
}

func (tsr *tsRetrieveQuery) AllTimeRange() storage.QueryCacheTSRetrieve {
	tsr.allRange = true
	return tsr
}

func (tsr *tsRetrieveQuery) WhereTimeRange(fromTS int64, toTS int64) storage.QueryCacheTSRetrieve {
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
			return tsr.modelExchange.bindRes(pk.String(), res)
		})
	}
	tsr.baseExchangeToSource()
	return tsr.baseErr()
}

func (tsr *tsRetrieveQuery) validateReq() {
	tsr.catcher.Exec(func() error {
		return tsRetrieveQueryReqValidator.Exec(tsr).Error()
	})
}

var tsRetrieveQueryReqValidator = validate.New[*tsRetrieveQuery]([]func(q *tsRetrieveQuery) error{
	validatePKProvided,
})

func validatePKProvided(q *tsRetrieveQuery) error {
	if (len(q.PKChain)) == 0 {
		return storage.Error{
			Type:    storage.ErrorTypeInvalidArgs,
			Message: "no PKC provided to ts retrieve query",
		}
	}
	return nil
}
