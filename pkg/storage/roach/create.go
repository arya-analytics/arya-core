package roach

import (
	"context"
	"github.com/uptrace/bun"
)

type Create struct {
	query *bun.InsertQuery
}

func NewCreate(c *bun.DB) *Create {
	q := c.NewInsert()
	return &Create{query: q}
}

func (c *Create) Model(model interface{}) *Create {
	c.query = c.query.Model(model)
	return c
}

func (c *Create) Exec(ctx context.Context) error {
	// TODO: Handle SQL Result Here
	_, err := c.query.Exec(ctx)
	return err
}
