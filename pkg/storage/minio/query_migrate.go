package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
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
		bucketExists, err := m.baseClient().BucketExists(ctx, wrap.Bucket())
		if err != nil {
			return m.baseHandleExecErr(err)
		}
		if !bucketExists {
			if mErr := m.baseClient().MakeBucket(ctx, wrap.Bucket(),
				minio.MakeBucketOptions{}); mErr != nil {
				return m.baseHandleExecErr(mErr)
			}
		} else {
			log.Warnf("Found bucket %s that already exists.", wrap.Bucket())
		}
	}
	return nil
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
