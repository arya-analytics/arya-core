package redis

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	log "github.com/sirupsen/logrus"
)

type TSQueryVariant int

const (
	TSQueryVariantSeries TSQueryVariant = iota
	TSQueryVariantSample
)

type tsCreateQuery struct {
	variant TSQueryVariant
	baseQuery
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
	tsc.catcher.Exec(func() error {
		wrapper := &TSModelWrapper{rfl: tsc.modelAdapter.Dest()}
		switch tsc.variant {
		case TSQueryVariantSample:
			return tsc.baseClient().TSCreateSamples(ctx, wrapper.Samples()...).Err()
		case TSQueryVariantSeries:
			for _, in := range wrapper.SeriesNames() {
				if err := tsc.baseClient().TSCreateSeries(ctx, in,
					timeseries.CreateOptions{}).Err(); err != nil {
					return err
				}
			}
		default:
			return storage.Error{
				Type:    storage.ErrTypeInvalidArgs,
				Message: "ts create queries require a variant selection",
			}
		}
		return nil
	})
	return tsc.baseErr()
}
