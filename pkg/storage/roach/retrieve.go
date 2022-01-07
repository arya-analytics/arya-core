package roach

import (
	"context"
	"github.com/uptrace/bun"
)

type Retrieve struct {
	query *bun.SelectQuery
}

func NewRetrieve(c *bun.DB) *Retrieve {
	q := c.NewSelect()
	return &Retrieve{query: q}
}

func (r *Retrieve) Model(model interface{}) *Retrieve {
	r.query = r.query.Model(model)
	return r
}

func (r *Retrieve) Where(query string, args ...interface{}) *Retrieve {
	r.query = r.query.Where(query, args...)
	return r
}

func (r *Retrieve) Exec(ctx context.Context) error {
	// TODO: Handle this error gracefully
	return r.query.Scan(ctx)
}
