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
	r.q = r.q.Model(r.baseModel(m).Pointer())
	return r
}

func (r *retrieveQuery) Where(query string, args ...interface{}) storage.MDRetrieveQuery {
	r.q = r.q.Where(query, args...)
	return r
}

func (r *retrieveQuery) WherePK(pk interface{}) storage.MDRetrieveQuery {
	return r.Where("PK = ?", pk)
}

func (r *retrieveQuery) WherePKs(pks interface{}) storage.MDRetrieveQuery {
	return r.Where("PK in (?)", bun.In(pks))
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	err := r.q.Scan(ctx)
	if err != nil {
		return r.baseHandleExecErr(err)
	}
	r.baseAdaptToSource()
	return err
}
