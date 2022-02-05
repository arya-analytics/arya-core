package storage_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	ctx   = context.Background()
	store = storage.New(storage.Config{
		MDEngine:     roach.New(mock.DriverPG{}),
		ObjectEngine: minio.New(mock.DriverMinio{}),
		CacheEngine:  redis.New(mock.DriverRedis{}),
	})
)

var _ = BeforeSuite(func() {
	err := store.NewMigrate().Exec(ctx)
	Expect(err).To(BeNil())
})

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}
