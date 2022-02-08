package mock

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
)

type Storage struct {
	storage.Storage
	DriverRoach *DriverRoach
	DriverRedis DriverRedis
	DriverMinio DriverMinio
}

func (s Storage) Stop() {
	s.DriverRoach.Stop()
}

func NewStorage() *Storage {
	driverRoach := &DriverRoach{}
	driverMinio := DriverMinio{}
	driverRedis := DriverRedis{}

	return &Storage{
		Storage: storage.New(storage.Config{
			EngineMD:     roach.New(driverRoach),
			EngineCache:  redis.New(driverRedis),
			EngineObject: minio.New(driverMinio),
		}),
		DriverRoach: driverRoach,
		DriverMinio: driverMinio,
		DriverRedis: driverRedis,
	}
}
