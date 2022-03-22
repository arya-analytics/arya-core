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

func (c *Create) GoExec(ctx context.Context, e chan error) {
	NewGoExecOpt(c.Pack(), e)
	go func() {
		err := c.Exec(ctx)
		if err != nil {
			e <- err
		}
	}()
}
