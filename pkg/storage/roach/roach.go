package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
)

// |||| CONFIG ||||

type TransactionLogLevel int

const (
	// TransactionLogLevelNone logs no queries.
	TransactionLogLevelNone TransactionLogLevel = iota
	// TransactionLogLevelErr logs failed queries.
	TransactionLogLevelErr
	// TransactionLogLevelAll logs all queries.
	TransactionLogLevelAll
)

type Config struct {
}

// |||| ENGINE ||||

// Engine opens connections and execute queries with a roach database.
// implements the storage.MDEngine interface.
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

// NewCreate opens a new createQuery query with the provided storage.Adapter.
func (e *Engine) NewCreate(a storage.Adapter) storage.MDCreateQuery {
	return newCreate(conn(a))
}

// NewRetrieve opens a new retrieveQuery query with the provided storage.Adapter.
func (e *Engine) NewRetrieve(a storage.Adapter) storage.MDRetrieveQuery {
	return newRetrieve(conn(a))
}

// NewUpdate opens a new updateQuery with the provided storage.Adapter.
func (e *Engine) NewUpdate(a storage.Adapter) storage.MDUpdateQuery {
	return newUpdate(conn(a))
}

// NewDelete opens a new deleteQuery with the provided storage.Adapter.
func (e *Engine) NewDelete(a storage.Adapter) storage.MDDeleteQuery {
	return newDelete(conn(a))
}

// NewMigrate opens a new migrateQuery with the provided storage.Adapter.
func (e *Engine) NewMigrate(a storage.Adapter) storage.MDMigrateQuery {
	return newMigrate(conn(a), e.driver)
}
