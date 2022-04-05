package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// TSCreate is for writing Queries that persist time series data to storage.
type TSCreate struct {
	query.Create
}

// NewTSCreate creates a new TSCreate query.
func NewTSCreate() *TSCreate {
	c := &TSCreate{}
	c.Base.Init(c)
	return c
}

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct, pointer to a slice, or pointer to a channel.
// The model must contain all necessary values and satisfy any relationships.
//
// If running the query using Stream, the model must be a pointer to a channel.
func (c *TSCreate) Model(m interface{}) *TSCreate {
	c.Base.Model(m)
	return c
}

// BindExec binds Execute that TSCreate will use to run the query.
// This method MUST be called before calling Exec.
func (c *TSCreate) BindExec(exec query.Execute) *TSCreate {
	c.Base.BindExec(exec)
	return c
}

// BindStream binds a Stream that TSCreate will use to run the query.
func (c *TSCreate) BindStream(stream *Stream) *TSCreate {
	BindStreamOpt(c.Pack(), stream)
	return c
}

// Stream starts the stream for writing values to query.Execute.
// The Stream returned will pipe errors encountered during value streaming.
// The error returned as the second argument represents errors encountered during query assembly.
//
// To close the stream, call a context.CancelFunc instead of closing the channel.
//
func (c *TSCreate) Stream(ctx context.Context) (*Stream, error) {
	o, ok := StreamOpt(c.Pack())
	if !ok {
		o = NewStreamOpt(ctx, c.Pack())
	}
	return o, c.Exec(ctx)
}
