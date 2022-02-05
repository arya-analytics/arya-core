package storage

import (
	"context"
)

// CreateQuery creates a new model in storage.
type CreateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newCreate(s *Storage) *CreateQuery {
	c := &CreateQuery{}
	c.baseInit(s)
	return c
}

/// |||| INTERFACE ||||

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct or a pointer to a slice.
// The model must contain all necessary values and satisfy any relationships.
func (c *CreateQuery) Model(m interface{}) *CreateQuery {
	c.baseBindModel(m)
	c.setMDQuery(c.mdQuery().Model(c.modelRfl.Pointer()))
	return c
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (c *CreateQuery) Exec(ctx context.Context) error {
	c.baseExec(func() error { return c.mdQuery().Exec(ctx) })
	mp := c.modelRfl.Pointer()
	if c.baseObjEngine().InCatalog(mp) {
		c.baseExec(func() error { return c.objQuery().Model(mp).Exec(ctx) })
	}
	return c.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (c *CreateQuery) mdQuery() MDCreateQuery {
	if c.baseMDQuery() == nil {
		c.setMDQuery(c.baseMDEngine().NewCreate(c.baseMDAdapter()))
	}
	return c.baseMDQuery().(MDCreateQuery)
}

func (c *CreateQuery) setMDQuery(q MDCreateQuery) {
	c.baseSetMDQuery(q)
}

// || OBJECT ||

func (c *CreateQuery) objQuery() ObjectCreateQuery {
	if c.baseObjQuery() == nil {
		c.setObjQuery(c.baseObjEngine().NewCreate(c.baseObjAdapter()))
	}
	return c.baseObjQuery().(ObjectCreateQuery)
}

func (c *CreateQuery) setObjQuery(q ObjectCreateQuery) {
	c.baseSetObjQuery(q)
}
