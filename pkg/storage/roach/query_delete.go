package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type queryDelete struct {
	queryBase
	q *bun.DeleteQuery
}

func newDelete(db *bun.DB) *queryDelete {
	q := &queryDelete{q: db.NewDelete()}
	q.baseInit()
	return q
}

func (q *queryDelete) WherePK(pk interface{}) storage.QueryMDDelete {
	return q.Where(pkEqualsSQL, pk)
}

func (q *queryDelete) WherePKs(pks interface{}) storage.QueryMDDelete {
	return q.Where(pkChainInSQL, bun.In(pks))
}

func (q *queryDelete) Where(query string, args ...interface{}) storage.QueryMDDelete {
	q.q = q.q.Where(query, args...)
	return q
}

func (q *queryDelete) Model(m interface{}) storage.QueryMDDelete {
	q.q = q.q.Model(q.baseModel(m).Pointer())
	return q
}

func (q *queryDelete) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.q.Exec(ctx)
		return err
	})
	return q.baseErr()
}
