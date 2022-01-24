package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type deleteQuery struct {
	baseQuery
	q *bun.DeleteQuery
}

func newDelete(db *bun.DB) *deleteQuery {
	r := &deleteQuery{q: db.NewDelete()}
	return r
}

func (d *deleteQuery) WherePK(pks interface{}) storage.MDDeleteQuery {
	return d.Where("PK = ?", pks)
}

func (d *deleteQuery) WherePKs(pks interface{}) storage.MDDeleteQuery {
	return d.Where("PK in (?)", bun.In(pks))
}

func (d *deleteQuery) Where(query string, args ...interface{}) storage.MDDeleteQuery {
	d.q = d.q.Where(query, args...)
	return d
}

func (d *deleteQuery) Model(m interface{}) storage.MDDeleteQuery {
	d.q = d.q.Model(d.baseModel(m).Pointer())
	return d
}

func (d *deleteQuery) Exec(ctx context.Context) error {
	_, err := d.q.Exec(ctx)
	return err
}
