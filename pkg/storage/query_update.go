package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
	q.setMDQuery(q.mdQuery().Model(model))
	return q
}

// WherePK queries the primary key of the model to be deleted.
func (q *QueryUpdate) WherePK(pk interface{}) *QueryUpdate {
	q.setMDQuery(q.mdQuery().WherePK(pk))
	return q
}

// Fields marks the fields that need to be updated. If this option isn't specified,
// will replace all fields.
//
// NOTE: When calling Bulk, order matters. Fields must be called before Bulk.
func (q *QueryUpdate) Fields(flds ...string) *QueryUpdate {
	q.setMDQuery(q.mdQuery().Fields(flds...))
	return q
}

// Bulk marks the update as a bulk update and allows for
// the update of multiple records.
func (q *QueryUpdate) Bulk() *QueryUpdate {
	q.setMDQuery(q.mdQuery().Bulk())
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

func (q *QueryUpdate) mdQuery() *query.Update {
	if q.baseMDQuery() == nil {
		q.setMDQuery(q.baseMDEngine().NewUpdate())
	}
	return q.baseMDQuery().(*query.Update)
}

func (q *QueryUpdate) setMDQuery(qmd *query.Update) {
	q.baseSetMDQuery(qmd)
}
