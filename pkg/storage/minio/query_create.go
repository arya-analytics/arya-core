package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var putObjectOpts = minio.PutObjectOptions{}

type createQuery struct {
	queryBase
}

func newCreate(client *minio.Client) *createQuery {
	c := &createQuery{}
	c.baseInit(client)
	return c
}

func (c *createQuery) Model(m interface{}) storage.QueryObjectCreate {
	c.baseModel(m)
	c.baseExchangeToDest()
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	for _, dv := range c.modelExchange.DataVals() {
		c.catcher.Exec(func() error {
			if dv.Data == nil {
				return storage.Error{
					Type:    storage.ErrorTypeInvalidArgs,
					Message: fmt.Sprintf("Minio data to write is nil! Model %s with id %s", c.modelExchange.Dest.Type(), dv.PK),
				}
			}
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
