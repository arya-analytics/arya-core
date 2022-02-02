package storage

import "context"

type UpdateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newUpdate(s *Storage) *UpdateQuery {
	u := &UpdateQuery{}
	u.baseInit(s)
	return u
}

func (u *UpdateQuery) Model(model interface{}) *UpdateQuery {
	u.setMDQuery(u.mdQuery().Model(model))
	return u
}

func (u *UpdateQuery) WherePK(pk interface{}) *UpdateQuery {
	u.setMDQuery(u.mdQuery().WherePK(pk))
	return u
}

func (u *UpdateQuery) Exec(ctx context.Context) error {
	u.catcher.Exec(func() error {
		return u.mdQuery().Exec(ctx)
	})
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
