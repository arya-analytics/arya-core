package mock

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
)

func NewStorage() *storage.Storage {
	return storage.New(storage.Config{
		MDEngine:     roach.New(DriverPG{}),
		CacheEngine:  redis.New(DriverRedis{}),
		ObjectEngine: minio.New(DriverMinio{}),
	})
}
