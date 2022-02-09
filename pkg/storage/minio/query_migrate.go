package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var makeBucketOpts = minio.MakeBucketOptions{}

type queryMigrate struct {
	queryBase
}

func newMigrate(client *minio.Client) *queryMigrate {
	q := &queryMigrate{}
	q.baseInit(client)
	return q
}

func (q *queryMigrate) Exec(ctx context.Context) error {
	for _, mod := range catalog() {
		me := newWrappedModelExchange(storage.NewModelExchange(mod, mod))
		q.catcher.Exec(func() error {
			bucketExists, err := q.baseClient().BucketExists(ctx, me.Bucket())
			if err != nil {
				return err
			}
			if !bucketExists {
				if mErr := q.baseClient().MakeBucket(ctx, me.Bucket(),
					makeBucketOpts); mErr != nil {
					return mErr
				}
			}
			return nil
		})
	}
	return q.baseErr()
}

func (q *queryMigrate) Verify(ctx context.Context) error {
	for _, mod := range catalog() {
		me := newWrappedModelExchange(storage.NewModelExchange(mod, mod))
		exists, err := q.baseClient().BucketExists(ctx, me.Bucket())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("bucket %s does not exist", err)
		}
	}
	return nil
}
