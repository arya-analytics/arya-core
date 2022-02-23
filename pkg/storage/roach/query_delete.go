package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type queryDelete struct {
	queryBase
	bunQ *bun.DeleteQuery
}

func newDelete(db *bun.DB) *queryDelete {
	q := &queryDelete{bunQ: db.NewDelete()}
	q.baseInit(db)
	return q
}

func (q *queryDelete) WherePK(pk interface{}) storage.QueryMDDelete {
	return q.Where(q.baseSQL().pk(), pk)
}

func (q *queryDelete) WherePKs(pks interface{}) storage.QueryMDDelete {
	return q.Where(q.baseSQL().pks(), bun.In(pks))
}

func (q *queryDelete) Where(query string, args ...interface{}) storage.QueryMDDelete {
	q.bunQ = q.bunQ.Where(query, args...)
	return q
}

func (q *queryDelete) Model(m interface{}) storage.QueryMDDelete {
	q.baseModel(m)
	q.bunQ = q.bunQ.Model(q.Dest().Pointer())
	return q
}

func (q *queryDelete) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.bunQ.Exec(ctx)
		return err
	})
	return q.baseErr()
}
