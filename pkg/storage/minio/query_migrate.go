package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/minio/minio-go/v7"
)

type migrateQuery struct {
	baseQuery
}

func newMigrate(client *minio.Client) *migrateQuery {
	m := &migrateQuery{}
	m.baseInit(client)
	return m
}

func (m *migrateQuery) Exec(ctx context.Context) error {
	for _, mod := range catalog() {
		wrap := &ModelWrapper{rfl: model.NewReflect(mod)}
		m.catcher.Exec(func() error {
			bucketExists, err := m.baseClient().BucketExists(ctx, wrap.Bucket())
			if err != nil {
				return err
			}
			if !bucketExists {
				if mErr := m.baseClient().MakeBucket(ctx, wrap.Bucket(),
					minio.MakeBucketOptions{}); mErr != nil {
					return mErr
				}
			}
			return nil
		})
	}
	return m.baseErr()
}

func (m *migrateQuery) Verify(ctx context.Context) error {
	for _, mod := range catalog() {
		wrap := &ModelWrapper{rfl: model.NewReflect(mod)}
		exists, err := m.baseClient().BucketExists(ctx, wrap.Bucket())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("bucket %s does not exist", err)
		}
	}
	return nil
}
