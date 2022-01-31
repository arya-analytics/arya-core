package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type createQuery struct {
	baseQuery
	modelRfl *model.Reflect
}

// |||| CONSTRUCTOR ||||

func newCreate(s *Storage) *createQuery {
	c := &createQuery{}
	c.baseInit(s)
	return c
}

/// |||| INTERFACE ||||

func (c *createQuery) Model(m interface{}) *createQuery {
	c.modelRfl = model.NewReflect(m)
	c.setMDQuery(c.mdQuery().Model(c.modelRfl.Pointer()))
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	c.catcher.Exec(func() error {
		return c.mdQuery().Exec(ctx)
	})
	if c.storage.cfg.objEngine().InCatalog(c.modelRfl.Pointer()) {
		c.catcher.Exec(func() error {
			return c.objQuery().Model(c.modelRfl.Pointer()).Exec(ctx)
		})
	}
	return c.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (c *createQuery) mdQuery() MDCreateQuery {
	if c.baseMDQuery() == nil {
		c.setMDQuery(c.baseMDEngine().NewCreate(c.baseMDAdapter()))
	}
	return c.baseMDQuery().(MDCreateQuery)
}

func (c *createQuery) setMDQuery(q MDCreateQuery) {
	c.baseSetMDQuery(q)
}

// || OBJECT ||

func (c *createQuery) objQuery() ObjectCreateQuery {
	if c.baseObjQuery() == nil {
		c.setObjQuery(c.baseObjEngine().NewCreate(c.baseObjAdapter()))
	}
	return c.baseObjQuery().(ObjectCreateQuery)
}

func (c *createQuery) setObjQuery(q ObjectCreateQuery) {
	c.baseSetObjQuery(q)
}
