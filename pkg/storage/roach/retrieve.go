package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type retrieve struct {
	base
	q *bun.SelectQuery
}

func newRetrieve(db *bun.DB) *retrieve {
	r := &retrieve{q: db.NewSelect()}
	return r
}

func (r *retrieve) Model(m interface{}) storage.MetaDataRetrieve {
	r.bindWrappers(m)
	r.q = r.q.Model(r.roachWrapper.Model())
	return r
}

func (r *retrieve) Where(query string, args ...interface{}) storage.MetaDataRetrieve {
	// TODO: look into if this is actually necessary
	r.q = r.q.Where(query, args...)
	return r
}

func (r *retrieve) WhereID(id interface{}) storage.MetaDataRetrieve {
	return r.Where("ID = ?", id)
}

func (r *retrieve) Exec(ctx context.Context) error {
	err := r.q.Scan(ctx)
	if err != nil {
		return err
	}
	err = r.storageWrapper.BindVals(r.roachWrapper.MapVals())
	return err
}
