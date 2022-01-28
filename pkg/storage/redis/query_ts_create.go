package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	log "github.com/sirupsen/logrus"
)

type TSQueryVariant int

const (
	TSQueryVariantSeries TSQueryVariant = iota + 1
	TSQueryVariantSample
)

type tsCreateQuery struct {
	tsBaseQuery
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
	log.SetReportCaller(true)
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
	w := tsc.tsBaseModelWrapper()
	tsc.catcher.Exec(func() error {
		return tsc.baseClient().TSCreateSamples(ctx, w.samples()...).Err()
	})
}

func (tsc *tsCreateQuery) execSeries(ctx context.Context) {
	w := tsc.tsBaseModelWrapper()
	for _, in := range w.seriesNames() {
		tsc.catcher.Exec(func() error {
			return tsc.baseClient().TSCreateSeries(ctx, in,
				timeseries.CreateOptions{}).Err()
		})
	}
}
