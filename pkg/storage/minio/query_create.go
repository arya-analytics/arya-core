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
	mw := c.baseModelWrapper()
	for _, dv := range mw.DataVals() {
		_, err := c.baseClient().PutObject(
			ctx,
			mw.Bucket(),
			dv.PK.String(),
			dv.Data,
			dv.Data.Size(),
			minio.PutObjectOptions{},
		)
		if err != nil {
			return c.baseHandleExecErr(err)
		}
	}
	return nil
}
