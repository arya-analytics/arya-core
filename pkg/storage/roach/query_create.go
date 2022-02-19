package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

// queryCreate implements storage.QueryMDCreate.
type queryCreate struct {
	queryBase
	q *bun.InsertQuery
}

func newCreate(db *bun.DB) *queryCreate {
	q := &queryCreate{q: db.NewInsert()}
	q.baseInit(db)
	return q
}

// Model implements storage.QueryMDCreate.
func (q *queryCreate) Model(m interface{}) storage.QueryMDCreate {
	q.baseModel(m)
	q.baseExchangeToDest()
	q.catcher.Exec(func() error {
		beforeInsertSetUUID(q.Dest())
		q.q = q.q.Model(q.Dest().Pointer())
		return nil
	})
	// We set base values, so we want to exchange back to source.
	q.baseExchangeToSource()
	return q
}

// Exec implements storage.QueryMDCreate.
func (q *queryCreate) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.q.Exec(ctx)
		return err
	})
	return q.baseErr()
}
