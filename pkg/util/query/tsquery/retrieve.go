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

func (r *Retrieve) WherePK(pk interface{}) *Retrieve {
	r.Where.WherePK(pk)
	return r
}

func (r *Retrieve) WherePKs(pks interface{}) *Retrieve {
	r.Where.WherePKs(pks)
	return r
}

func (r *Retrieve) AllTime() *Retrieve {
	NewTimeRangeOpt(r.Pack(), telem.AllTime())
	return r

}

func (r *Retrieve) WhereTimeRange(tr telem.TimeRange) {
	NewTimeRangeOpt(r.Pack(), tr)
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
