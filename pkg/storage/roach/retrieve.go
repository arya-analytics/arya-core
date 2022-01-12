package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type Retrieve struct {
	q *bun.SelectQuery
	model interface{}
	rModel interface{}
}

func newRetrieve(db *bun.DB) *Retrieve {
	r := &Retrieve{q: db.NewSelect()}
	return r
}

func (r *Retrieve) Model(m interface{}) storage.MetaDataRetrieve {
	r.model = m
	rm := roachModel(m)
	r.rModel = rm
	r.q = r.q.Model(rm)
	return r
}

func (r *Retrieve) Where(query string, args ...interface{}) storage.MetaDataRetrieve {
	r.q = r.q.Where(query, args...)
	return r
}

func (r *Retrieve) Exec(ctx context.Context) error {
	err := r.q.Scan(ctx)
	return err
}
