package storage

import "context"

type createQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newCreate(s *Storage) *createQuery {
	c := &createQuery{}
	c.baseInit(s)
	return c
}

/// ||| INT

func (c *createQuery) Model(model interface{}) *createQuery {
	c.baseSetMDQuery(c.mdQuery().Model(model))
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	return c.mdQuery().Exec(ctx)
}

// || META DATA QUERY BINDING ||

func (c *createQuery) mdQuery() MDCreateQuery {
	if c.baseMDQuery() == nil {
		c.setMDQuery(c.mdEngine.NewCreate(c.baseMDAdapter()))
	}
	return c.baseMDQuery().(MDCreateQuery)
}

func (c *createQuery) setMDQuery(q MDCreateQuery) {
	c.baseSetMDQuery(q)
}
