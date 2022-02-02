package storage

import "context"

type DeleteQuery struct {
	baseQuery
}

func newDelete(s *Storage) *DeleteQuery {
	d := &DeleteQuery{}
	d.baseInit(s)
	return d
}

// |||| INTERFACE ||||

func (d *DeleteQuery) Model(model interface{}) *DeleteQuery {
	d.setMDQuery(d.mdQuery().Model(model))
	return d
}

func (d *DeleteQuery) WherePK(pk interface{}) *DeleteQuery {
	d.setMDQuery(d.mdQuery().WherePK(pk))
	return d
}

func (d *DeleteQuery) WherePKs(pks interface{}) *DeleteQuery {
	d.setMDQuery(d.mdQuery().WherePKs(pks))
	return d
}

func (d *DeleteQuery) Exec(ctx context.Context) error {
	d.catcher.Exec(func() error {
		return d.mdQuery().Exec(ctx)
	})
	return d.baseErr()
}

// |||| QUERY BINDING ||||

func (d *DeleteQuery) mdQuery() MDDeleteQuery {
	if d.baseMDQuery() == nil {
		d.setMDQuery(d.baseMDEngine().NewDelete(d.baseMDAdapter()))
	}
	return d.baseMDQuery().(MDDeleteQuery)
}

func (d *DeleteQuery) setMDQuery(q MDDeleteQuery) {
	d.baseSetMDQuery(q)
}
