package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type updateQuery struct {
	baseQuery
	q *bun.UpdateQuery
}

func newUpdate(db *bun.DB) *updateQuery {
	r := &updateQuery{q: db.NewUpdate()}
	r.baseInit()
	return r
}

func (u *updateQuery) Model(m interface{}) storage.MDUpdateQuery {
	rm := u.baseModel(m)
	u.baseAdaptToDest()
	u.q = u.q.Model(rm.Pointer())
	return u
}

func (u *updateQuery) WherePK(pk interface{}) storage.MDUpdateQuery {
	return u.Where("ID = ?", pk)
}

func (u *updateQuery) Where(query string, args ...interface{}) storage.MDUpdateQuery {
	u.q = u.q.Where(query, args...)
	return u
}

func (u *updateQuery) Exec(ctx context.Context) error {
	u.catcher.Exec(func() error {
		_, err := u.q.Exec(ctx)
		return err
	})
	return u.baseErr()
}
