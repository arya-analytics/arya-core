package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

// createQuery implements storage.MDCreateQuery.
type createQuery struct {
	baseQuery
	q *bun.InsertQuery
}

func newCreate(db *bun.DB) *createQuery {
	r := &createQuery{q: db.NewInsert()}
	r.baseInit()
	return r
}

// Model implements storage.MDCreateQuery.
func (c *createQuery) Model(m interface{}) storage.MDCreateQuery {
	rm := c.baseModel(m)
	c.baseExchangeToDest()
	c.catcher.Exec(func() error {
		beforeInsertSetUUID(rm)
		c.q = c.q.Model(rm.Pointer())
		return nil
	})
	// We set default values, so we want to exchange back to source.
	c.baseExchangeToSource()
	return c
}

// Exec implements storage.MDCreateQuery.
func (c *createQuery) Exec(ctx context.Context) error {
	c.baseExec(func() error {
		_, err := c.q.Exec(ctx)
		return err
	})
	return c.baseErr()
}
