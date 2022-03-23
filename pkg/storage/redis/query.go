package redis

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/go-redis/redis/v8"
)

type base struct {
	client *timeseries.Client
	exc    *exchange
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
	tsc.exc.ToDest()
	if tsc.qt == queryVariantSeries {
		err = tsc.execSeries(ctx, p)
	} else {
		err = tsc.execSample(ctx, p)
	}
	tsc.exc.ToSource()
	return newErrorConvert().Exec(err)
}

func (tsr *tsRetrieve) exec(ctx context.Context, p *query.Pack) (err error) {
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
		if err := tsr.exc.bindRes(pk.String(), res); err != nil {
			return err
		}
	}
	tsr.exc.ToSource()
	return nil
}

// |||| OPT CONVERTERS ||||

func (tsc *tsCreate) convertOpts(p *query.Pack) {
	internal.OptConverters{tsc.model, tsc.variant}.Exec(p)
}

func (tsr *tsRetrieve) convertOpts(p *query.Pack) {
	internal.OptConverters{tsr.model, tsr.pk, tsr.timeRange}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) {
	b.exc = wrapExchange(model.NewExchange(p.Model(), catalog().New(p.Model())))
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
	tr, _ := tsquery.TimeRangeOpt(p)
	tsr.tr = tr
}

// |||| CUSTOM CREATE ||||

func (tsc *tsCreate) execSeries(ctx context.Context, p *query.Pack) error {
	for _, sn := range tsc.exc.seriesNames() {
		if err := tsc.base.client.TSCreateSeries(ctx, sn, timeseries.CreateOptions{}).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (tsc *tsCreate) execSample(ctx context.Context, p *query.Pack) error {
	return tsc.base.client.TSCreateSamples(ctx, tsc.exc.samples()...).Err()
}

type queryVariant int

const (
	queryVariantSeries queryVariant = iota + 1
	queryVariantSample
)

func (tsc *tsCreate) variant(p *query.Pack) {
	bt, ok := tsc.exc.Dest().StructTagChain().RetrieveBase()
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