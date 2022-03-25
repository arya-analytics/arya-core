package tsquery

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Create struct {
	query.Create
}

func NewCreate() *Create {
	c := &Create{}
	c.Base.Init(c)
	return c
}

func (c *Create) Model(m interface{}) *Create {
	c.Base.Model(m)
	return c
}

func (c *Create) BindExec(exec query.Execute) *Create {
	c.Base.BindExec(exec)
	return c
}

func (c *Create) GoExec(ctx context.Context) GoExecOpt {
	o := NewGoExecOpt(c.Pack())
	go func() {
		if err := c.Exec(ctx); err != nil {
			o.Errors <- err
		}
	}()
	return o
}
