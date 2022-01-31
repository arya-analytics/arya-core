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

func (d *deleteQuery) WherePK(pk interface{}) *deleteQuery {
	d.setMDQuery(d.mdQuery().WherePK(pk))
	return d
}

func (d *deleteQuery) WherePKs(pks interface{}) *deleteQuery {
	d.setMDQuery(d.mdQuery().WherePKs(pks))
	return d
}

func (d *deleteQuery) Exec(ctx context.Context) error {
	d.catcher.Exec(func() error {
		return d.mdQuery().Exec(ctx)
	})
	return d.baseErr()
}

// |||| QUERY BINDING ||||

func (d *deleteQuery) mdQuery() MDDeleteQuery {
	if d.baseMDQuery() == nil {
		d.setMDQuery(d.baseMDEngine().NewDelete(d.baseMDAdapter()))
	}
	return d.baseMDQuery().(MDDeleteQuery)
}

func (d *deleteQuery) setMDQuery(q MDDeleteQuery) {
	d.baseSetMDQuery(q)
}
