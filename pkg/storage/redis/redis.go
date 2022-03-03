package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
)

// |||| ENGINE ||||

type Engine struct {
	pool   *storage.Pool
	driver Driver
}

func New(driver Driver, pool *storage.Pool) *Engine {
	return &Engine{driver: driver, pool: pool}
}

func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.driver)
}

func (e *Engine) client() *timeseries.Client {
	return conn(e.pool.Retrieve(e))
}

func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) NewTSRetrieve() storage.QueryCacheTSRetrieve {
	return newTSRetrieve(e.client())
}

func (e *Engine) NewTSCreate() storage.QueryCacheTSCreate {
	return newTSCreate(e.client())
}
