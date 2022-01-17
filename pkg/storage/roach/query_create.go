package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

type createQuery struct {
	baseQuery
	q *bun.InsertQuery
}

func newCreate(db *bun.DB) *createQuery {
	r := &createQuery{q: db.NewInsert()}
	return r
}

func (c *createQuery) Model(m interface{}) storage.MDCreateQuery {
	rm := c.baseModel(m)
	c.baseBindErr(createValidator.Exec(m))
	c.baseAdaptToDest()
	c.q = c.q.Model(rm)
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	if c.baseCheckErr() {
		return c.baseErr()
	}
	_, err := c.q.Exec(ctx)
	return c.baseHandleExecErr(err)
}
