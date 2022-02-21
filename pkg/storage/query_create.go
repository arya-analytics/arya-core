package storage

import (
	"context"
)

// QueryCreate creates a new model in storage.
type QueryCreate struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newCreate(s Storage) *QueryCreate {
	c := &QueryCreate{}
	c.baseInit(s, s.config().Hooks.BeforeCreate)
	return c
}

/// |||| INTERFACE ||||

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct or a pointer to a slice.
// The model must contain all necessary values and satisfy any relationships.
func (c *QueryCreate) Model(m interface{}) *QueryCreate {
	c.baseBindModel(m)
	c.baseRunHook()
	c.setMDQuery(c.mdQuery().Model(c.modelRfl.Pointer()))
	return c
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (c *QueryCreate) Exec(ctx context.Context) error {
	c.baseExec(func() error { return c.mdQuery().Exec(ctx) })
	mp := c.modelRfl.Pointer()
	if c.baseObjEngine().InCatalog(mp) {
		c.baseExec(func() error { return c.objQuery().Model(mp).Exec(ctx) })
	}
	return c.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (c *QueryCreate) mdQuery() QueryMDCreate {
	if c.baseMDQuery() == nil {
		c.setMDQuery(c.baseMDEngine().NewCreate(c.baseMDAdapter()))
	}
	return c.baseMDQuery().(QueryMDCreate)
}

func (c *QueryCreate) setMDQuery(q QueryMDCreate) {
	c.baseSetMDQuery(q)
}

// || OBJECT ||

func (c *QueryCreate) objQuery() QueryObjectCreate {
	if c.baseObjQuery() == nil {
		c.setObjQuery(c.baseObjEngine().NewCreate(c.baseObjAdapter()))
	}
	return c.baseObjQuery().(QueryObjectCreate)
}

func (c *QueryCreate) setObjQuery(q QueryObjectCreate) {
	c.baseSetObjQuery(q)
}
