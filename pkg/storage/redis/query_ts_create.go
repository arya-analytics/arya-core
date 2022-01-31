package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
)

type TSQueryVariant int

const (
	TSQueryVariantSeries TSQueryVariant = iota + 1
	TSQueryVariantSample
)

type tsCreateQuery struct {
	baseQuery
	variant TSQueryVariant
}

func newTSCreate(client *timeseries.Client) *tsCreateQuery {
	tsc := &tsCreateQuery{}
	tsc.baseInit(client)
	return tsc
}

func (tsc *tsCreateQuery) Series() storage.CacheTSCreateQuery {
	tsc.variant = TSQueryVariantSeries
	return tsc
}

func (tsc *tsCreateQuery) Sample() storage.CacheTSCreateQuery {
	tsc.variant = TSQueryVariantSample
	return tsc
}

func (tsc *tsCreateQuery) Model(m interface{}) storage.CacheTSCreateQuery {
	tsc.baseModel(m)
	tsc.baseAdaptToDest()
	return tsc
}

func (tsc *tsCreateQuery) Exec(ctx context.Context) error {
	switch tsc.variant {
	case TSQueryVariantSample:
		tsc.execSample(ctx)
	case TSQueryVariantSeries:
		tsc.execSeries(ctx)
	default:
		return storage.Error{
			Type:    storage.ErrTypeInvalidArgs,
			Message: "ts create queries require a variant selection",
		}
	}
	return tsc.baseErr()
}

func (tsc *tsCreateQuery) execSample(ctx context.Context) {
	tsc.catcher.Exec(func() error {
		return tsc.baseClient().TSCreateSamples(ctx, tsc.modelAdapter.samples()...).Err()
	})
}

func (tsc *tsCreateQuery) execSeries(ctx context.Context) {
	for _, in := range tsc.modelAdapter.seriesNames() {
		tsc.catcher.Exec(func() error {
			return tsc.baseClient().TSCreateSeries(ctx, in,
				timeseries.CreateOptions{}).Err()
		})
	}
}
