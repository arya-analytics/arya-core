package storage

import "context"

type deleteQuery struct {
	baseQuery
}

func newDelete(s *Storage) *deleteQuery {
	d := &deleteQuery{}
	d.baseInit(s)
	return d
}

// |||| INTERFACE ||||

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

// |||| QUERY BINDING ||||

func (d *deleteQuery) mdQuery() MDDeleteQuery {
	if d.baseMDQuery() == nil {
		d.setMDQuery(d.mdEngine.NewDelete(d.baseMDAdapter()))
	}
	return d.baseMDQuery().(MDDeleteQuery)
}

func (d *deleteQuery) setMDQuery(q MDDeleteQuery) {
	d.baseSetMDQuery(q)
}
