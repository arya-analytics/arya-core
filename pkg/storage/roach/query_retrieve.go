package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type queryRetrieve struct {
	queryBase
	q *bun.SelectQuery
}

func newRetrieve(db *bun.DB) *queryRetrieve {
	q := &queryRetrieve{q: db.NewSelect()}
	q.baseInit(db)
	return q
}

func (q *queryRetrieve) Model(m interface{}) storage.QueryMDRetrieve {
	q.baseModel(m)
	q.q = q.q.Model(q.Dest().Pointer())
	return q
}

func (q *queryRetrieve) Where(query string, args ...interface{}) storage.QueryMDRetrieve {
	q.q = q.q.Where(query, args...)
	return q
}

func (q *queryRetrieve) WherePK(pk interface{}) storage.QueryMDRetrieve {
	return q.Where(q.baseSQL().pk(), pk)
}

func (q *queryRetrieve) WherePKs(pks interface{}) storage.QueryMDRetrieve {
	return q.Where(q.baseSQL().pks(), bun.In(pks))
}

func (q *queryRetrieve) Relation(rel string) storage.QueryMDRetrieve {
	q.q = q.q.Relation(rel)
	return q
}

func (q *queryRetrieve) Count(ctx context.Context) (count int, err error) {
	q.baseExec(func() error {
		count, err = q.q.Count(ctx)
		return err
	})
	return count, q.baseErr()
}

func (q *queryRetrieve) Exec(ctx context.Context) error {
	q.baseExec(func() error { return q.q.Scan(ctx) })
	q.baseExchangeToSource()
	return q.baseErr()
}
