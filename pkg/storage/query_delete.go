package storage

import "context"

// DeleteQuery deletes a model or set of models depending on the parameters passed.
// DeleteQuery requires that WherePK or WherePKs is called,
// and will panic upon DeleteQuery.Exec if it isn't.
//
// DeleteQuery should not be instantiated directly,
// and should instead should be opened using Storage.NewDelete().
type DeleteQuery struct {
	baseQuery
}

func newDelete(s *Storage) *DeleteQuery {
	d := &DeleteQuery{}
	d.baseInit(s)
	return d
}

// |||| INTERFACE ||||

// Model sets the model to be deleted. model must be passed as a pointer.
// DeleteQuery only uses the model for location lookups,
// and doesn't actually do anything  with the value.
// You're good to provide a nil pointer as long as it has the right type.
func (d *DeleteQuery) Model(model interface{}) *DeleteQuery {
	d.baseBindModel(model)
	d.setMDQuery(d.mdQuery().Model(model))
	return d
}

// WherePK queries by the primary of the model to be deleted.
func (d *DeleteQuery) WherePK(pk interface{}) *DeleteQuery {
	d.setMDQuery(d.mdQuery().WherePK(pk))
	return d
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (d *DeleteQuery) WherePKs(pks interface{}) *DeleteQuery {
	d.setMDQuery(d.mdQuery().WherePKs(pks))
	return d
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (d *DeleteQuery) Exec(ctx context.Context) error {
	d.catcher.Exec(func() error { return d.mdQuery().Exec(ctx) })
	if d.baseObjEngine().InCatalog(d.modelRfl.Pointer()) {
		d.catcher.Exec(func() error {
			return d.objQuery().Model(d.modelRfl.Pointer()).WherePKs(d.modelRfl.
				PKChain().Raw(),
			).Exec(ctx)
		})
	}
	return d.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (d *DeleteQuery) mdQuery() MDDeleteQuery {
	if d.baseMDQuery() == nil {
		d.setMDQuery(d.baseMDEngine().NewDelete(d.baseMDAdapter()))
	}

	return d.baseMDQuery().(MDDeleteQuery)
}

func (d *DeleteQuery) setMDQuery(q MDDeleteQuery) {
	d.baseSetMDQuery(q)
}

// || OBJECT ||

func (d *DeleteQuery) objQuery() ObjectDeleteQuery {
	if d.baseObjQuery() == nil {
		d.setObjQuery(d.baseObjEngine().NewDelete(d.baseObjAdapter()))
	}
	return d.baseObjQuery().(ObjectDeleteQuery)
}

func (d *DeleteQuery) setObjQuery(q ObjectDeleteQuery) {
	d.baseSetObjQuery(q)
}
