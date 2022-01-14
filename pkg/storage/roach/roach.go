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

// engine opens connections and execute queries with a roach database.
// implements the storage.MDEngine interface.
type engine struct {
	cfg Config
}

func New(cfg Config) *engine {
	return &engine{cfg}
}

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *engine) NewAdapter() storage.Adapter {
	return newAdapter(e.cfg)
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

// NewRetrieve opens a new retrieveQuery query with the provided storage.Adapter.
func (e *engine) NewRetrieve(a storage.Adapter) storage.MDRetrieveQuery {
	return newRetrieve(conn(a))
}

// NewCreate opens a new createQuery query with the provided storage.Adapter.
func (e *engine) NewCreate(a storage.Adapter) storage.MDCreateQuery {
	return newCreate(conn(a))
}

// NewDelete opens a new deleteQuery with the provided storage.Adapter;
func (e *engine) NewDelete(a storage.Adapter) storage.MDDeleteQuery {
	return newDelete(conn(a))
}

// NewMigrate opens a new migrateQuery with the provided storage.Adapter;
func (e *engine) NewMigrate(a storage.Adapter) storage.MigrateQuery {
	return newMigrate(conn(a), e.cfg.Driver)
}
