package storage

import "context"

type createQuery struct {
	baseQuery
	_mdQuery MetaDataCreate
}

func newCreate(s *Storage) *createQuery {
	c := &createQuery{}
	c.init(s)
	return c
}

//TODO: use generics here
func (c *createQuery) mdQuery() MetaDataCreate {
	if c._mdQuery == nil {
		c._mdQuery = c.mdEngine.NewCreate(c.storage.adapter(EngineRoleMetaData))
	}
	return c._mdQuery
}

func (c *createQuery) setMdQuery(q MetaDataCreate) {
	c._mdQuery = q
}

func (c *createQuery) Model(model interface{}) *createQuery {
	c.setMdQuery(c.mdQuery().Model(model))
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	return c.mdQuery().Exec(ctx)
}
