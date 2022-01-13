package storage

import "context"

type create struct {
	base
	_mdQuery MetaDataCreate
}

func newCreate(s *Storage) *create {
	c := &create{}
	c.init(s)
	return c
}

//TODO: use generics here
func (c *create) mdQuery() MetaDataCreate {
	if c._mdQuery == nil {
		c._mdQuery = c.mdEngine.NewCreate(c.storage.adapter(EngineRoleMetaData))
	}
	return c._mdQuery
}

func (c *create) setMdQuery(q MetaDataCreate) MetaDataCreate {
	c._mdQuery = q
	return c._mdQuery
}

func (c *create) Model(model interface{}) *create {
	c.setMdQuery(c.mdQuery().Model(model))
	return c
}

func (c *create) Exec(ctx context.Context) error {
	return c.mdQuery().Exec(ctx)
}
