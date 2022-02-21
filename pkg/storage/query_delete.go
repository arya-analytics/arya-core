package storage

import "context"

// QueryDelete deletes a model or set of models depending on the parameters passed.
// QueryDelete requires that WherePK or WherePKs is called,
// and will panic upon QueryDelete.Exec if it isn't.
//
// QueryDelete should not be instantiated directly,
// and should instead should be opened using Storage.NewDelete().
type QueryDelete struct {
	queryBase
}

func newDelete(s Storage) *QueryDelete {
	q := &QueryDelete{}
	q.baseInit(s, q)
	return q
}

// |||| INTERFACE ||||

// Model sets the model to be deleted. model must be passed as a pointer.
// QueryDelete only uses the model for location lookups,
// and doesn't actually do anything  with the value.
// You're good to provide a nil pointer as long as it has the right type.
func (q *QueryDelete) Model(model interface{}) *QueryDelete {
	q.baseBindModel(model)
	q.setMDQuery(q.mdQuery().Model(model))
	return q
}

// WherePK queries by the primary of the model to be deleted.
func (q *QueryDelete) WherePK(pk interface{}) *QueryDelete {
	q.setMDQuery(q.mdQuery().WherePK(pk))
	return q
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (q *QueryDelete) WherePKs(pks interface{}) *QueryDelete {
	q.setMDQuery(q.mdQuery().WherePKs(pks))
	return q
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (q *QueryDelete) Exec(ctx context.Context) error {
	q.baseRunBeforeHooks(ctx)
	q.baseExec(func() error { return q.mdQuery().Exec(ctx) })
	if q.baseObjEngine().InCatalog(q.modelRfl.Pointer()) {
		q.baseExec(func() error {
			return q.objQuery().Model(q.modelRfl.Pointer()).WherePKs(q.modelRfl.PKChain().Raw()).Exec(ctx)
		})
	}
	q.baseRunAfterHooks(ctx)
	return q.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (q *QueryDelete) mdQuery() QueryMDDelete {
	if q.baseMDQuery() == nil {
		q.setMDQuery(q.baseMDEngine().NewDelete(q.baseMDAdapter()))
	}

	return q.baseMDQuery().(QueryMDDelete)
}

func (q *QueryDelete) setMDQuery(qmd QueryMDDelete) {
	q.baseSetMDQuery(qmd)
}

// || OBJECT ||

func (q *QueryDelete) objQuery() QueryObjectDelete {
	if q.baseObjQuery() == nil {
		q.setObjQuery(q.baseObjEngine().NewDelete(q.baseObjAdapter()))
	}
	return q.baseObjQuery().(QueryObjectDelete)
}

func (q *QueryDelete) setObjQuery(qob QueryObjectDelete) {
	q.baseSetObjQuery(qob)
}
