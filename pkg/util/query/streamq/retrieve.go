package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

type TSRetrieve struct {
	query.Retrieve
}

func NewTSRetrieve() *TSRetrieve {
	r := &TSRetrieve{}
	r.base.Init(r)
	return r
}

func (r *TSRetrieve) Model(m interface{}) *TSRetrieve {
	r.base.Model(m)
	return r
}

func (r *TSRetrieve) WherePKs(pks interface{}) *TSRetrieve {
	r.Retrieve.WherePKs(pks)
	return r
}

func (r *TSRetrieve) WherePK(pk interface{}) *TSRetrieve {
	r.Retrieve.WherePK(pk)
	return r
}

func (r *TSRetrieve) AllTime() *TSRetrieve {
	return r.WhereTimeRange(telem.AllTime())
}

func (r *TSRetrieve) WhereTimeRange(tr telem.TimeRange) *TSRetrieve {
	NewTimeRangeOpt(r.Pack(), tr)
	return r
}

func (r *TSRetrieve) BindExec(exec query.Execute) *TSRetrieve {
	r.base.BindExec(exec)
	return r
}

func (r *TSRetrieve) BindStream(stream *Stream) *TSRetrieve {
	BindStreamOpt(r.Pack(), stream)
	return r
}

func (r *TSRetrieve) Stream(ctx context.Context) (*Stream, error) {
	o, ok := StreamOpt(r.Pack())
	if !ok {
		o = NewStreamOpt(ctx, r.Pack())
	}
	return o, r.Exec(ctx)
}

const timeRangeOptKey query.OptKey = "tsRange"

func NewTimeRangeOpt(p *query.Pack, tr telem.TimeRange) {
	p.SetOpt(timeRangeOptKey, tr)
}

func TimeRangeOpt(p *query.Pack) (telem.TimeRange, bool) {
	opt, ok := p.RetrieveOpt(timeRangeOptKey)
	if !ok {
		return telem.TimeRange{}, false
	}
	return opt.(telem.TimeRange), true
}
