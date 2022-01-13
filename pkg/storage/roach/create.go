package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type create struct {
	base
	q *bun.InsertQuery
}

func newCreate(db *bun.DB) *create {
	r := &create{q: db.NewInsert()}
	return r
}

func (c *create) Model(m interface{}) storage.MetaDataCreate {
	c.bindWrappers(m)
	if err := c.roachWrapper.BindVals(c.storageWrapper.MapVals()); err != nil {
		log.Fatalln(err)
	}
	c.q = c.q.Model(c.roachWrapper.Model())
	return c
}

func (c *create) Exec(ctx context.Context) error {
	_, err := c.q.Exec(ctx)
	return err
}
