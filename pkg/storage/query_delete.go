package storage

import "context"

type deleteQuery struct {
	base
	_mdQuery MetaDataDelete
}

func newDelete(s *Storage) *deleteQuery {
	d := &deleteQuery{}
	d.init(s)
	return d
}

func (d *deleteQuery) mdQuery() MetaDataDelete {
	if d._mdQuery == nil {
		d._mdQuery = d.mdEngine.NewDelete(d.storage.adapter(EngineRoleMetaData))
	}
	return d._mdQuery
}

func (d *deleteQuery) setMDQuery(q MetaDataDelete) {
	d._mdQuery = q
}

func (d *deleteQuery) Model(model interface{}) *deleteQuery {
	d.setMDQuery(d.mdQuery().Model(model))
	return d
}

func (d *deleteQuery) WhereID(id interface{}) *deleteQuery {
	d.setMDQuery(d.mdQuery().WhereID(id))
	return d
}

func (d *deleteQuery) Exec(ctx context.Context) error {
	return d.mdQuery().Exec(ctx)
}
