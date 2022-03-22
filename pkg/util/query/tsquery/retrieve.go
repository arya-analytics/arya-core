package tsquery

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

type Retrieve struct {
	query.Retrieve
}

func NewRetrieve() *Retrieve {
	r := &Retrieve{}
	r.Base.Init(r)
	return r
}

func (r *Retrieve) Model(m interface{}) *Retrieve {
	r.Base.Model(m)
	return r
}

func (r *Retrieve) AllTime() *Retrieve {
	return r.WhereTimeRange(telem.AllTime())
}

func (r *Retrieve) WhereTimeRange(tr telem.TimeRange) *Retrieve {
	NewTimeRangeOpt(r.Pack(), tr)
	return r
}

func (r *Retrieve) BindExec(exec query.Execute) *Retrieve {
	r.Base.BindExec(exec)
	return r
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
