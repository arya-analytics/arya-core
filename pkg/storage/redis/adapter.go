package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type adapter struct {
	id     uuid.UUID
	client *timeseries.Client
	cfg    Config
}

func newAdapter(cfg Config) *adapter {
	a := &adapter{
		id:  uuid.New(),
		cfg: cfg,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func conn(a storage.Adapter) *timeseries.Client {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind redis adapter.")
	}
	return ra.conn()
}

func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) open() {
	switch a.cfg.Driver {
	case DriverRedisTS:
		a.client = connectToRedis(a.cfg)
	}
}

func (a *adapter) conn() *timeseries.Client {
	return a.client
}

func redisConfig(cfg Config) *redis.Options {
	return &redis.Options{
		Addr:     cfg.addr(),
		DB:       cfg.Database,
		Password: cfg.Password,
	}
}

func connectToRedis(cfg Config) *timeseries.Client {
	return timeseries.NewWrap(redis.NewClient(redisConfig(cfg)))
}
