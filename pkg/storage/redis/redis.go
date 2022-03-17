package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	log "github.com/sirupsen/logrus"
)

// |||| ENGINE ||||

type Engine struct {
	pool   *storage.Pool
	driver Driver
}

func New(driver Driver, pool *storage.Pool) *Engine {
	return &Engine{driver: driver, pool: pool}
}

func (e *Engine) NewAdapter() (storage.Adapter, error) {
	return newAdapter(e.driver)
}

func (e *Engine) client() (*timeseries.Client, error) {
	a, err := e.pool.Retrieve(e)
	if err != nil {
		return nil, err
	}
	return conn(a), nil
}

func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) NewTSRetrieve() storage.QueryCacheTSRetrieve {
	c, err := e.client()
	if err != nil {
		log.Fatalln(err)
	}
	return newTSRetrieve(c)
}

func (e *Engine) NewTSCreate() storage.QueryCacheTSCreate {
	c, err := e.client()
	if err != nil {
		log.Fatalln(err)
	}
	return newTSCreate(c)
}
