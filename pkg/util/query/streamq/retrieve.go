package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

// TSRetrieve is for writing Queries that retrieve time series data from a data store.
type TSRetrieve struct {
	query.Retrieve
}

// NewTSRetrieve opens a new TSRetrieve query.
func NewTSRetrieve() *TSRetrieve {
	r := &TSRetrieve{}
	r.Base.Init(r)
	return r
}

// Model sets the model to bind the results into. model must be passed as a pointer.
// If you're expecting multiple return values,
// pass a pointer to a slice. If you're expecting one return value,
// pass a struct. NOTE: If a struct is passed, and multiple values are returned,
// the struct is assigned to the value of the first result.
//
// If running the query using Stream, the model must be a pointer to a channel.
func (r *TSRetrieve) Model(m interface{}) *TSRetrieve {
	r.Base.Model(m)
	return r
}

// WherePK queries by the primary key of the model to be retrieved.
func (r *TSRetrieve) WherePK(pk interface{}) *TSRetrieve {
	r.Retrieve.WherePK(pk)
	return r
}

// WherePKs queries by a set of primary keys of models to be retrieved.
func (r *TSRetrieve) WherePKs(pks interface{}) *TSRetrieve {
	r.Retrieve.WherePKs(pks)
	return r
}

// AllTime queries across the entire time span for the model.
func (r *TSRetrieve) AllTime() *TSRetrieve {
	return r.WhereTimeRange(telem.AllTime())
}

// WhereTimeRange queries across a specific time range for the model.
func (r *TSRetrieve) WhereTimeRange(tr telem.TimeRange) *TSRetrieve {
	NewTimeRangeOpt(r.Pack(), tr)
	return r
}

// BindExec binds Execute that TSRetrieve will use to run the query.
// This method MUST be called before calling Exec.
func (r *TSRetrieve) BindExec(exec query.Execute) *TSRetrieve {
	r.Base.BindExec(exec)
	return r
}

// BindStream binds a Stream that TSRetrieve will use to run the query.
func (r *TSRetrieve) BindStream(stream *Stream) *TSRetrieve {
	BindStreamOpt(r.Pack(), stream)
	return r
}

// Stream starts the stream that sends values to the passed model.
// The Stream returned will pipe errors encountered during value streaming.
// The error returned as the second  argument represents errors encountered during query assembly.
//
// To close the stream, call a context.CancelFunc instead of closing the channel.
// DO NOT CLOSE THE CHANNEL. This will cause the stream to panic.
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

func TimeRangeOpt(p *query.Pack, opts ...query.OptRetrieveOpt) (telem.TimeRange, bool) {
	opt, ok := p.RetrieveOpt(timeRangeOptKey, opts...)
	if !ok {
		return telem.TimeRange{}, false
	}
	return opt.(telem.TimeRange), true
}
