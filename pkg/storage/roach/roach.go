package roach

import (
	"crypto/tls"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
)

// |||| CONFIG ||||

type Driver int

const (
	DriverPG Driver = iota
	DriverSQLite
)

type Config struct {
	// Username for the database. Does not need to be specified if using DriverSQLite.
	Username string
	// Password for the database. Does not need to be specified if using DriverSQLite.
	Password string
	// Host IP for the database. Does not need to be specified if using DriverSQLite.
	Host string
	// Port to connect to at Host. Does not need to be specified if using DriverSQLite.
	Port int
	// Database to connect to. Does not need to be specified if using DriverSQLite.
	Database string
	// Whether to open a TLS connection or not.
	// Does not need to be specified if using DriverSQLite.
	UseTLS bool
	// Driver is the connection driver used for the roach data store.
	// Current options are:
	// DriverPG which connects via the Postgres wire protocol.
	// DriverSQLite which uses an in memory SQLite database
	Driver Driver
}

func (e Config) addr() string {
	return fmt.Sprintf("%s:%v", e.Host, e.Port)
}

func (e Config) tls() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: e.UseTLS,
	}
}

// |||| ENGINE ||||

// Engine opens connections and execute queries with a roach database.
// implements the storage.MetaDataEngine interface.
type Engine struct {
	cfg Config
}

func New(cfg Config) *Engine {
	return &Engine{cfg}
}

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.cfg)
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

// NewRetrieve opens a new retrieveQuery query with the provided storage.Adapter.
func (e *Engine) NewRetrieve(a storage.Adapter) storage.MetaDataRetrieve {
	return newRetrieve(conn(a))
}

// NewCreate opens a new createQuery query with the provided storage.Adapter.
func (e *Engine) NewCreate(a storage.Adapter) storage.MetaDataCreate {
	return newCreate(conn(a))
}

// NewDelete opens a new deleteQuery with the provided storage.Adapter;
func (e *Engine) NewDelete(a storage.Adapter) storage.MetaDataDelete {
	return newDelete(conn(a))
}

// NewMigrate opens a new migrate with the provided storage.Adapter;
func (e *Engine) NewMigrate(a storage.Adapter) storage.Migrate {
	return newMigrate(conn(a), e.cfg.Driver)
}
