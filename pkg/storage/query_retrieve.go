package storage

import (
	"context"
)

// RetrieveQuery retrieves a model or set of models depending on the parameters passed.
// RetrieveQuery requires that WherePK or WherePKs is called,
// and will panic upon RetrieveQuery.Exec if it is not.
//
// RetrieveQuery should not be instantiated directly,
// and should instead be opened using Storage.NewRetrieve().
type RetrieveQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s *Storage) *RetrieveQuery {
	r := &RetrieveQuery{}
	r.baseInit(s)
	return r
}

// |||| INTERFACE ||||

// Model sets the model to bind the results into. model must be passed as a pointer.
// If you're expecting multiple return values,
// pass a pointer to a slice. If you're expecting one return value,
// pass a struct. NOTE: If a struct is passed, and multiple values are returned,
// we assign the value of the first result.
func (r *RetrieveQuery) Model(m interface{}) *RetrieveQuery {
	r.baseBindModel(m)
	r.setMDQuery(r.mdQuery().Model(m))
	return r
}

// WherePK queries by the primary of the model to be deleted.
func (r *RetrieveQuery) WherePK(pk interface{}) *RetrieveQuery {
	r.setMDQuery(r.mdQuery().WherePK(pk))
	return r
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (r *RetrieveQuery) WherePKs(pks interface{}) *RetrieveQuery {
	r.setMDQuery(r.mdQuery().WherePKs(pks))
	return r
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (r *RetrieveQuery) Exec(ctx context.Context) error {
	r.baseExec(func() error { return r.mdQuery().Exec(ctx) })
	if r.baseObjEngine().InCatalog(r.modelRfl.Pointer()) {
		r.baseExec(func() error {
			return r.objQuery().Model(r.modelRfl.Pointer()).WherePKs(r.modelRfl.PKChain().Raw()).
				Exec(ctx)
		})
	}
	return r.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (r *RetrieveQuery) mdQuery() MDRetrieveQuery {
	if r.baseMDQuery() == nil {
		r.setMDQuery(r.baseMDEngine().NewRetrieve(r.baseMDAdapter()))
	}
	return r.baseMDQuery().(MDRetrieveQuery)
}

func (r *RetrieveQuery) setMDQuery(q MDRetrieveQuery) {
	r.baseSetMDQuery(q)
}

// || OBJECT ||

func (r *RetrieveQuery) objQuery() ObjectRetrieveQuery {
	if r.baseObjQuery() == nil {
		r.setObjQuery(r.baseObjEngine().NewRetrieve(r.baseObjAdapter()))
	}
	return r.baseObjQuery().(ObjectRetrieveQuery)
}

func (r *RetrieveQuery) setObjQuery(q ObjectRetrieveQuery) {
	r.baseSetObjQuery(q)
}
