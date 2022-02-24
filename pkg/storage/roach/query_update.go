package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type queryUpdate struct {
	queryBase
	bunQ *bun.UpdateQuery
}

func newUpdate(db *bun.DB) *queryUpdate {
	q := &queryUpdate{bunQ: db.NewUpdate()}
	q.baseInit(db)
	return q
}

func (q *queryUpdate) Model(m interface{}) storage.QueryMDUpdate {
	q.baseModel(m)
	q.baseExchangeToDest()
	q.bunQ = q.bunQ.Model(q.baseDest().Pointer())
	return q
}

func (q *queryUpdate) Bulk() storage.QueryMDUpdate {
	q.bunQ = q.bunQ.Bulk().OmitZero()
	return q
}

func (q *queryUpdate) Fields(flds ...string) storage.QueryMDUpdate {
	q.bunQ = q.bunQ.Column(q.baseSQL().fieldNames(flds...)...)
	return q
}

func (q *queryUpdate) WherePK(pk interface{}) storage.QueryMDUpdate {
	q.bunQ = q.bunQ.Where(q.baseSQL().pk(), pk)
	return q
}

func (q *queryUpdate) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.bunQ.Exec(ctx)
		return err
	})
	return q.baseErr()
}
