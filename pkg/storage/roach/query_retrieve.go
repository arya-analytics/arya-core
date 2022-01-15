package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type retrieveQuery struct {
	baseQuery
	q *bun.SelectQuery
}

func newRetrieve(db *bun.DB) *retrieveQuery {
	r := &retrieveQuery{q: db.NewSelect()}
	return r
}

func (r *retrieveQuery) Model(m interface{}) storage.MDRetrieveQuery {
	r.q = r.q.Model(r.baseModel(m))
	return r
}

func (r *retrieveQuery) Where(query string, args ...interface{}) storage.MDRetrieveQuery {
	r.q = r.q.Where(query, args...)
	return r
}

func (r *retrieveQuery) WhereID(id interface{}) storage.MDRetrieveQuery {
	return r.Where("ID = ?", id)
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	err := r.q.Scan(ctx)
	if err != nil {
		return err
	}
	r.baseAdaptToSource()
	return err
}
