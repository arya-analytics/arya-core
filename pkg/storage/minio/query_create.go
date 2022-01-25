package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

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
	c.baseAdaptToDest()
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	for _, dv := range c.baseModelWrapper().DataVals() {
		c.catcher.Exec(func() error {
			_, err := c.baseClient().PutObject(
				ctx,
				c.Bucket(),
				dv.PK.String(),
				dv.Data,
				dv.Data.Size(),
				minio.PutObjectOptions{},
			)
			return err
		})
	}
	return c.baseErr()
}
