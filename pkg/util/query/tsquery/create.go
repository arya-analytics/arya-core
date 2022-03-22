package tsquery

import "github.com/arya-analytics/aryacore/pkg/util/query"

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
