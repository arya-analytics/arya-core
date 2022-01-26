package redists

import (
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
)

type Adapter struct {
	id     uuid.UUID
	client *redistimeseries.Client
	cfg    Config
}

func newAdapter(cfg Config) *Adapter {
	a := &Adapter{
		id:  uuid.New(),
		cfg: cfg,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*Adapter, bool) {
	ra, ok := a.(*Adapter)
	return ra, ok
}

func conn(a storage.Adapter) *redistimeseries.Client {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind redists adapter.")
	}
	return ra.conn()
}

func (a *Adapter) ID() uuid.UUID {
	return a.id
}

func (a *Adapter) open() {
	switch a.cfg.Driver {
	case DriverRedisTS:
		a.client = connectToRedisTs(a.cfg)
	}
}

func (a *Adapter) conn() *redistimeseries.Client {
	return a.client
}

func connectToRedisTs(cfg Config) *redistimeseries.Client {
	return redistimeseries.NewClient(cfg.addr(), cfg.Database, cfg.authString())
}
