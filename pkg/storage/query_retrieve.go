package storage

import (
	"context"
)

// QueryRetrieve retrieves a model or set of models depending on the parameters passed.
// QueryRetrieve requires that WherePK or WherePKs is called,
// and will panic upon QueryRetrieve.Exec if it is not.
//
// QueryRetrieve should not be instantiated directly,
// and should instead be opened using Storage.NewRetrieve().
type QueryRetrieve struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s Storage) *QueryRetrieve {
	q := &QueryRetrieve{}
	q.baseInit(s)
	return q
}

// |||| INTERFACE ||||

// Model sets the model to bind the results into. model must be passed as a pointer.
// If you're expecting multiple return values,
// pass a pointer to a slice. If you're expecting one return value,
// pass a struct. NOTE: If a struct is passed, and multiple values are returned,
// the struct is assigned to the value of the first result.
func (q *QueryRetrieve) Model(m interface{}) *QueryRetrieve {
	q.baseBindModel(m)
	q.setMDQuery(q.mdQuery().Model(m))
	return q
}

// WherePK queries by the primary of the model to be deleted.
func (q *QueryRetrieve) WherePK(pk interface{}) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().WherePK(pk))
	return q
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (q *QueryRetrieve) WherePKs(pks interface{}) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().WherePKs(pks))
	return q
}

func (q *QueryRetrieve) Relation(rel string, fields ...string) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().Relation(rel, fields...))
	return q
}

func (q *QueryRetrieve) Field(fields ...string) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().Field(fields...))
	return q
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (q *QueryRetrieve) Exec(ctx context.Context) error {
	q.baseExec(func() error { return q.mdQuery().Exec(ctx) })
	mp := q.modelRfl.Pointer()
	if q.baseObjEngine().InCatalog(mp) {
		pks := q.modelRfl.PKChain().Raw()
		q.baseExec(func() error { return q.objQuery().Model(mp).WherePKs(pks).Exec(ctx) })
	}
	return q.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (q *QueryRetrieve) mdQuery() QueryMDRetrieve {
	if q.baseMDQuery() == nil {
		q.setMDQuery(q.baseMDEngine().NewRetrieve(q.baseMDAdapter()))
	}
	return q.baseMDQuery().(QueryMDRetrieve)
}

func (q *QueryRetrieve) setMDQuery(qmd QueryMDRetrieve) {
	q.baseSetMDQuery(qmd)
}

// || OBJECT ||

func (q *QueryRetrieve) objQuery() QueryObjectRetrieve {
	if q.baseObjQuery() == nil {
		q.setObjQuery(q.baseObjEngine().NewRetrieve(q.baseObjAdapter()))
	}
	return q.baseObjQuery().(QueryObjectRetrieve)
}

func (q *QueryRetrieve) setObjQuery(qob QueryObjectRetrieve) {
	q.baseSetObjQuery(qob)
}
