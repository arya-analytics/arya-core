package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
)

// |||| ENGINE ||||

type Engine struct {
	driver Driver
}

func New(driver Driver) *Engine {
	return &Engine{driver}
}

func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.driver)
}

func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) InCatalog(m interface{}) bool {
	return catalog().Contains(m)
}

func (e *Engine) NewTSRetrieve(a storage.Adapter) storage.CacheTSRetrieveQuery {
	return newTSRetrieve(conn(a))
}

func (e *Engine) NewTSCreate(a storage.Adapter) storage.CacheTSCreateQuery {
	return newTSCreate(conn(a))
}
