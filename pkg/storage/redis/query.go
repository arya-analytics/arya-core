package redis

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/go-redis/redis/v8"
)

type base struct {
	client       *timeseries.Client
	wrappedModel *reflectRedis
}

type tsCreate struct {
	base
	qt queryVariant
}

func newTSCreate(client *timeseries.Client) *tsCreate {
	return &tsCreate{base: base{client: client}}
}

type tsRetrieve struct {
	base
	pkc model.PKChain
	tr  telem.TimeRange
}

func newTSRetrieve(client *timeseries.Client) *tsRetrieve {
	return &tsRetrieve{base: base{client: client}}
}

// |||| EXEC ||||

func (tsc *tsCreate) exec(ctx context.Context, p *query.Pack) (err error) {
	tsc.convertOpts(p)
	if tsc.qt == queryVariantSeries {
		err = tsc.execSeries(ctx, p)
	} else {
		err = tsc.execSample(ctx, p)
	}
	return err
}

func (tsr *tsRetrieve) exec(ctx context.Context, p *query.Pack) error {
	tsr.convertOpts(p)
	for _, pk := range tsr.pkc {
		var cmd *redis.Cmd
		if tsr.tr.IsZero() {
			cmd = tsr.client.TSGet(ctx, pk.String())
		} else {
			cmd = tsr.base.client.TSGetRange(ctx, pk.String(), tsr.tr)
		}
		res, err := cmd.Result()
		if err != nil {
			return err
		}
		if bErr := tsr.wrappedModel.bindRes(pk.String(), res); bErr != nil {
			return bErr
		}
	}
	return nil
}

// |||| OPT CONVERTERS ||||

func (tsc *tsCreate) convertOpts(p *query.Pack) {
	query.OptConvertChain{tsc.model, tsc.variant}.Exec(p)
}

func (tsr *tsRetrieve) convertOpts(p *query.Pack) {
	query.OptConvertChain{tsr.model, tsr.pk, tsr.timeRange}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) {
	b.wrappedModel = wrapReflect(p.Model())
}

// |||| PK ||||

func (tsr *tsRetrieve) pk(p *query.Pack) {
	pkc, ok := query.PKOpt(p)
	if !ok {
		panic("tsRetrieve queries require a primary key!")
	}
	tsr.pkc = pkc
}

// |||| TIME RANGE ||||

func (tsr *tsRetrieve) timeRange(p *query.Pack) {
	tr, _ := streamq.RetrieveTimeRangeOpt(p)
	tsr.tr = tr
}

// |||| CUSTOM CREATE ||||

func (tsc *tsCreate) execSeries(ctx context.Context, p *query.Pack) error {
	for _, sn := range tsc.wrappedModel.seriesNames() {
		if err := tsc.base.client.TSCreateSeries(ctx, sn, timeseries.CreateOptions{}).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (tsc *tsCreate) execSample(ctx context.Context, p *query.Pack) error {
	c := tsc.base.client.TSCreateSamples(ctx, tsc.wrappedModel.samples()...)
	return c.Err()
}

type queryVariant int

const (
	queryVariantSeries queryVariant = iota + 1
	queryVariantSample
)

func (tsc *tsCreate) variant(p *query.Pack) {
	bt, ok := tsc.wrappedModel.StructTagChain().RetrieveBase()
	if !ok {
		panic(fmt.Sprintf("model %s does not specify a base tag", p.Model()))
	}
	v, ok := bt.Retrieve(model.TagCat, model.RoleKey)
	if !ok {
		panic(fmt.Sprintf("model %s did not specify a role", p.Model()))
	}
	if v == "tsSample" {
		tsc.qt = queryVariantSample
	} else if v == "tsSeries" {
		tsc.qt = queryVariantSeries
	} else {
		panic(fmt.Sprintf("model %s provided to query did not have correct role. Must by tsSample or tsSeries. Received %s", p.Model(), v))
	}
}
