package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var putObjectOpts = minio.PutObjectOptions{}

type createQuery struct {
	baseQuery
}

func newCreate(client *minio.Client) *createQuery {
	c := &createQuery{}
	c.baseInit(client)
	return c
}

func (c *createQuery) Model(m interface{}) storage.ObjectCreateQuery {
	c.baseModel(m)
	c.baseExchangeToDest()
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	for _, dv := range c.modelExchange.DataVals() {
		c.catcher.Exec(func() error {
			_, err := c.baseClient().PutObject(
				ctx,
				c.baseBucket(),
				dv.PK.String(),
				dv.Data,
				dv.Data.Size(),
				putObjectOpts,
			)
			return err
		})
	}
	return c.baseErr()
}
