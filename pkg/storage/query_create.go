package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// QueryCreate creates a new model in storage.
type QueryCreate struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newCreate(s Storage) *QueryCreate {
	c := &QueryCreate{}
	c.baseInit(s, c)
	return c
}

/// |||| INTERFACE ||||

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct or a pointer to a slice.
// The model must contain all necessary values and satisfy any relationships.
func (c *QueryCreate) Model(m interface{}) *QueryCreate {
	c.baseBindModel(m)
	c.setMDQuery(c.mdQuery().Model(c.modelRfl.Pointer()))
	return c
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (c *QueryCreate) Exec(ctx context.Context) error {
	c.baseRunBeforeHooks(ctx)
	c.baseExec(func() error { return c.mdQuery().Exec(ctx) })
	mp := c.modelRfl.Pointer()
	if c.baseObjEngine().ShouldHandle(mp) {
		c.baseExec(func() error { return c.objQuery().Model(mp).Exec(ctx) })
	}
	c.baseRunAfterHooks(ctx)
	return c.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (c *QueryCreate) mdQuery() *query.Create {
	if c.baseMDQuery() == nil {
		c.setMDQuery(c.baseMDEngine().NewCreate())
	}
	return c.baseMDQuery().(*query.Create)
}

func (c *QueryCreate) setMDQuery(q *query.Create) {
	c.baseSetMDQuery(q)
}

// || OBJECT ||

func (c *QueryCreate) objQuery() *query.Create {
	if c.baseObjQuery() == nil {
		c.setObjQuery(c.baseObjEngine().NewCreate())
	}
	return c.baseObjQuery().(*query.Create)
}

func (c *QueryCreate) setObjQuery(q *query.Create) {
	c.baseSetObjQuery(q)
}
