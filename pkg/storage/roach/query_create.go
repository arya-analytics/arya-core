package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
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

func (c *createQuery) Model(m interface{}) storage.MetaDataCreate {
	c.bindWrappers(m)
	if err := c.roachWrapper.BindVals(c.storageWrapper.MapVals()); err != nil {
		log.Fatalln(err)
	}
	c.q = c.q.Model(c.roachWrapper.Model())
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	_, err := c.q.Exec(ctx)
	return err
}
