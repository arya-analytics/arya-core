package storage

import (
	"context"
	"fmt"
)

// QueryUpdate updates a model in storage.
// QueryUpdate requires that WherePK is called, and will panic upon QueryUpdate.
// Exec if it is not.
type QueryUpdate struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newUpdate(s Storage) *QueryUpdate {
	q := &QueryUpdate{}
	q.baseInit(s, q)
	return q
}

// Model sets the model to update.
// The model MUST be a single struct and MUST be a pointer.
// NOTE: This query currently updates ALL values of the model. Not just defined ones.
func (q *QueryUpdate) Model(model interface{}) *QueryUpdate {
	q.baseBindModel(model)
	if !q.modelRfl.IsStruct() {
		panic(fmt.Sprintf("received a non struct model of type %T. "+
			"model must be a struct", q.modelRfl.Type()))
	}
	q.setMDQuery(q.mdQuery().Model(model))
	return q
}

// WherePK queries the primary key of the model to be deleted.
func (q *QueryUpdate) WherePK(pk interface{}) *QueryUpdate {
	q.setMDQuery(q.mdQuery().WherePK(pk))
	return q
}

// Exec execute the query with the provided context. Returns a storage.Error.
func (q *QueryUpdate) Exec(ctx context.Context) error {
	q.baseRunBeforeHooks(ctx)
	q.baseExec(func() error { return q.mdQuery().Exec(ctx) })
	q.baseRunAfterHooks(ctx)
	return q.baseErr()
}

// |||| QUERY BINDING ||||

func (q *QueryUpdate) mdQuery() QueryMDUpdate {
	if q.baseMDQuery() == nil {
		q.setMDQuery(q.baseMDEngine().NewUpdate(q.baseMDAdapter()))
	}
	return q.baseMDQuery().(QueryMDUpdate)
}

func (q *QueryUpdate) setMDQuery(qmd QueryMDUpdate) {
	q.baseSetMDQuery(qmd)
}
