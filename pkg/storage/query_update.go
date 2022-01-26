package storage

import "context"

type updateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newUpdate(s *Storage) *updateQuery {
	u := &updateQuery{}
	u.baseInit(s)
	return u
}

func (u *updateQuery) Model(model interface{}) *updateQuery {
	u.setMDQuery(u.mdQuery().Model(model))
	return u
}

func (u *updateQuery) WherePK(pk interface{}) *updateQuery {
	u.setMDQuery(u.mdQuery().WherePK(pk))
	return u
}

func (u *updateQuery) Exec(ctx context.Context) error {
	u.catcher.Exec(func() error {
		return u.mdQuery().Exec(ctx)
	})
	return u.baseErr()
}

// |||| QUERY BINDING ||||

func (u *updateQuery) mdQuery() MDUpdateQuery {
	if u.baseMDQuery() == nil {
		u.setMDQuery(u.mdEngine.NewUpdate(u.baseMDAdapter()))
	}
	return u.baseMDQuery().(MDUpdateQuery)
}

func (u *updateQuery) setMDQuery(q MDUpdateQuery) {
	u.baseSetMDQuery(q)
}
