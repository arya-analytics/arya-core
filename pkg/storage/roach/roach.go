package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/uptrace/bun"
)

// |||| CONFIG ||||

type Driver interface {
	Connect() (*bun.DB, error)
}

// |||| ENGINE ||||

// Engine opens connections and execute queries with a roach database.
// implements the storage.EngineMD interface.
type Engine struct {
	driver Driver
}

func New(driver Driver) *Engine {
	return &Engine{driver}
}

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.driver)
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) InCatalog(m interface{}) bool {
	return catalog().Contains(m)
}

// NewCreate opens a new queryCreate query with the provided storage.Adapter.
func (e *Engine) NewCreate(a storage.Adapter) storage.QueryMDCreate {
	return newCreate(conn(a))
}

// NewRetrieve opens a new queryRetrieve query with the provided storage.Adapter.
func (e *Engine) NewRetrieve(a storage.Adapter) storage.QueryMDRetrieve {
	return newRetrieve(conn(a))
}

// NewUpdate opens a new queryUpdate with the provided storage.Adapter.
func (e *Engine) NewUpdate(a storage.Adapter) storage.QueryMDUpdate {
	return newUpdate(conn(a))
}

// NewDelete opens a new queryDelete with the provided storage.Adapter.
func (e *Engine) NewDelete(a storage.Adapter) storage.QueryMDDelete {
	return newDelete(conn(a))
}

// NewMigrate opens a new queryMigrate with the provided storage.Adapter.
func (e *Engine) NewMigrate(a storage.Adapter) storage.QueryMDMigrate {
	return newMigrate(conn(a), e.driver)
}

func (e *Engine) Tasks(a storage.Adapter) *tasks.Scheduler {
	return newTaskScheduler(conn(a))
}
