package storage

import (
	"context"
	"fmt"
)

// UpdateQuery updates a model in storage.
// UpdateQuery requires that WherePK is called, and will panic upon UpdateQuery.
// Exec if it is not.
type UpdateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newUpdate(s *Storage) *UpdateQuery {
	u := &UpdateQuery{}
	u.baseInit(s)
	return u
}

// Model sets the model to update.
// The model MUST be a single struct and MUST be a pointer.
// NOTE: This query currently updates ALL values of the model. Not just defined ones.
func (u *UpdateQuery) Model(model interface{}) *UpdateQuery {
	u.baseBindModel(model)
	if !u.modelRfl.IsStruct() {
		panic(fmt.Sprintf("received a non struct model of type %T. "+
			"model must be a struct", u.modelRfl.Type()))
	}
	u.setMDQuery(u.mdQuery().Model(model))
	return u
}

// WherePK queries the primary key of the model to be deleted.
func (u *UpdateQuery) WherePK(pk interface{}) *UpdateQuery {
	u.setMDQuery(u.mdQuery().WherePK(pk))
	return u
}

// Exec execute the query with the provided context. Returns a storage.Error.
func (u *UpdateQuery) Exec(ctx context.Context) error {
	u.catcher.Exec(func() error { return u.mdQuery().Exec(ctx) })
	return u.baseErr()
}

// |||| QUERY BINDING ||||

func (u *UpdateQuery) mdQuery() MDUpdateQuery {
	if u.baseMDQuery() == nil {
		u.setMDQuery(u.baseMDEngine().NewUpdate(u.baseMDAdapter()))
	}
	return u.baseMDQuery().(MDUpdateQuery)
}

func (u *UpdateQuery) setMDQuery(q MDUpdateQuery) {
	u.baseSetMDQuery(q)
}
