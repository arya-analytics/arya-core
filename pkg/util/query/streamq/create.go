package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type TSCreate struct {
	query.Create
}

func NewTSCreate() *TSCreate {
	c := &TSCreate{}
	c.base.Init(c)
	return c
}

func (c *TSCreate) Model(m interface{}) *TSCreate {
	c.base.Model(m)
	return c
}

func (c *TSCreate) BindExec(exec query.Execute) *TSCreate {
	c.base.BindExec(exec)
	return c
}

func (c *TSCreate) BindStream(stream *Stream) *TSCreate {
	BindStreamOpt(c.Pack(), stream)
	return c
}

func (c *TSCreate) Stream(ctx context.Context) (*Stream, error) {
	o, ok := StreamOpt(c.Pack())
	if !ok {
		o = NewStreamOpt(ctx, c.Pack())
	}
	return o, c.Exec(ctx)
}
