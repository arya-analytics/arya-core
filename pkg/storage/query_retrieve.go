package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// QueryRetrieve retrieves a model or set of models depending on the parameters passed.
// QueryRetrieve requires that WherePK or WherePKs is called,
// and will panic upon QueryRetrieve.Exec if it is not.
//
// QueryRetrieve should not be instantiated directly,
// and should instead be opened using Storage.NewRetrieve().
type QueryRetrieve struct {
	queryBase
	_flds []string
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s Storage) *QueryRetrieve {
	q := &QueryRetrieve{}
	q.baseInit(s, q)
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

// WhereFields queries by a set of key value pairs where the key represents a field name
// and the value represent a value to match with.
func (q *QueryRetrieve) WhereFields(flds query.WhereFields) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().WhereFields(flds))
	return q
}

// Relation retrieves the specified fields from the relation.
func (q *QueryRetrieve) Relation(rel string, flds ...string) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().Relation(rel, flds...))
	return q
}

// Calculate executes a calculation on the specified field. It binds the calculation
// into the argument 'into'. See Calculate for available calculations.
func (q *QueryRetrieve) Calculate(c query.Calc, fldName string, into interface{}) *QueryRetrieve {
	q.setMDQuery(q.mdQuery().Calc(c, fldName, into))
	return q
}

// Fields retrieves only the fields specified.
func (q *QueryRetrieve) Fields(flds ...string) *QueryRetrieve {
	q._flds = flds
	q.setMDQuery(q.mdQuery().Fields(flds...))
	return q
}

func (q *QueryRetrieve) fields() (flds []string) {
	shouldHandleOmitFields := []string{"ID"}
out:
	for _, fld := range q._flds {
		for _, omit := range shouldHandleOmitFields {
			if fld == omit {
				continue out
			}
		}
		flds = append(flds, fld)
	}
	return flds
}

// Exec executes the query with the provided context. Returns a storage.Error.
func (q *QueryRetrieve) Exec(ctx context.Context) error {
	q.baseRunBeforeHooks(ctx)
	q.baseExec(func() error { return q.mdQuery().Exec(ctx) })
	mp := q.modelRfl.Pointer()
	if q.baseObjEngine().ShouldHandle(mp, q.fields()...) {
		pks := q.modelRfl.PKChain().Raw()
		q.baseExec(func() error { return q.objQuery().Model(mp).WherePKs(pks).Exec(ctx) })
	}
	q.baseRunAfterHooks(ctx)
	return q.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (q *QueryRetrieve) mdQuery() *query.Retrieve {
	if q.baseMDQuery() == nil {
		q.setMDQuery(q.baseMDEngine().NewRetrieve())
	}
	return q.baseMDQuery().(*query.Retrieve)
}

func (q *QueryRetrieve) setMDQuery(qmd *query.Retrieve) {
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
