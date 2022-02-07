package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type queryUpdate struct {
	queryBase
	q *bun.UpdateQuery
}

func newUpdate(db *bun.DB) *queryUpdate {
	q := &queryUpdate{q: db.NewUpdate()}
	q.baseInit()
	return q
}

func (q *queryUpdate) Model(m interface{}) storage.QueryMDUpdate {
	rm := q.baseModel(m)
	q.baseExchangeToDest()
	q.q = q.q.Model(rm.Pointer())
	return q
}

func (q *queryUpdate) WherePK(pk interface{}) storage.QueryMDUpdate {
	return q.Where(pkEqualsSQL, pk)
}

func (q *queryUpdate) Where(query string, args ...interface{}) storage.QueryMDUpdate {
	q.q = q.q.Where(query, args...)
	return q
}

func (q *queryUpdate) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.q.Exec(ctx)
		return err
	})
	return q.baseErr()
}
