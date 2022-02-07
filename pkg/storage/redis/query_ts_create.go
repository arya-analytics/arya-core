package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
)

type tsQueryVariant int

const (
	tsQueryVariantSeries tsQueryVariant = iota + 1
	tsQueryVariantSample
)

type tsCreateQuery struct {
	queryBase
	variant tsQueryVariant
}

func newTSCreate(client *timeseries.Client) *tsCreateQuery {
	tsc := &tsCreateQuery{}
	tsc.baseInit(client)
	return tsc
}

func (tsc *tsCreateQuery) Series() storage.QueryCacheTSCreate {
	tsc.variant = tsQueryVariantSeries
	return tsc
}

func (tsc *tsCreateQuery) Sample() storage.QueryCacheTSCreate {
	tsc.variant = tsQueryVariantSample
	return tsc
}

func (tsc *tsCreateQuery) Model(m interface{}) storage.QueryCacheTSCreate {
	tsc.baseModel(m)
	tsc.baseExchangeToDest()
	return tsc
}

func (tsc *tsCreateQuery) Exec(ctx context.Context) error {
	switch tsc.variant {
	case tsQueryVariantSample:
		tsc.execSample(ctx)
	case tsQueryVariantSeries:
		tsc.execSeries(ctx)
	default:
		return storage.Error{
			Type:    storage.ErrorTypeInvalidArgs,
			Message: "ts create queries require a variant selection",
		}
	}
	return tsc.baseErr()
}

func (tsc *tsCreateQuery) execSample(ctx context.Context) {
	tsc.catcher.Exec(func() error {
		return tsc.baseClient().TSCreateSamples(ctx, tsc.modelExchange.samples()...).Err()
	})
}

func (tsc *tsCreateQuery) execSeries(ctx context.Context) {
	for _, in := range tsc.modelExchange.seriesNames() {
		tsc.catcher.Exec(func() error {
			return tsc.baseClient().TSCreateSeries(ctx, in,
				timeseries.CreateOptions{}).Err()
		})
	}
}
