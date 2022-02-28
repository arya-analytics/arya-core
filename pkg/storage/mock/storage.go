package mock

import (
	"github.com/arya-analytics/aryacore/pkg/models"
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

func (s *Storage) Stop() {
	s.DriverRoach.Stop()
}

type storageOpts struct {
	Verbose bool
}

type StorageOpt func(so *storageOpts)

func WithVerbose() StorageOpt {
	return func(so *storageOpts) {
		so.Verbose = true
	}

}

func NewStorage(opts ...StorageOpt) *Storage {
	so := &storageOpts{}
	for _, opt := range opts {
		opt(so)
	}
	driverRoach := NewDriverRoach(false, so.Verbose)
	driverMinio := DriverMinio{}
	driverRedis := DriverRedis{}

	s := &Storage{
		Storage: storage.New(storage.Config{
			EngineMD:     roach.New(driverRoach),
			EngineCache:  redis.New(driverRedis),
			EngineObject: minio.New(driverMinio),
		}),
		DriverRoach: driverRoach,
		DriverMinio: driverMinio,
		DriverRedis: driverRedis,
	}
	models.BindHooks(s)
	return s
}
