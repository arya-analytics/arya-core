package roach

import (
	"crypto/tls"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
)

// |||| CONFIG ||||

type TransactionLogLevel int

const (
	// TransactionLogLevelNone logs no queries
	TransactionLogLevelNone TransactionLogLevel = iota
	// TransactionLogLevelErr logs failed queries
	TransactionLogLevelErr
	// TransactionLogLevelAll logs all queries1
	TransactionLogLevelAll
)

type Driver int

const (
	// DriverPG connects via the Postgres wire protocol.
	DriverPG Driver = iota
	// DriverSQLite which uses an in memory SQLite database
	DriverSQLite
)

type Config struct {
	// DSN is a connection string for the database. If specified,
	// all other fields except for Driver can be left blank.
	DSN string
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
	// Options are:
	Driver Driver
	// TransactionLogLevel is the log level for executed SQL queries
	TransactionLogLevel TransactionLogLevel
}

func (c Config) addr() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

func (c Config) tls() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: !c.UseTLS,
	}
}

// |||| ENGINE ||||

// Engine opens connections and execute queries with a roach database.
// implements the storage.MDEngine interface.
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

// IsAdapter checks if the provided Adapter is a roach Adapter.
func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) InCatalog(m interface{}) bool {
	return catalog().Contains(m)
}

// NewRetrieve opens a new retrieveQuery query with the provided storage.Adapter.
func (e *Engine) NewRetrieve(a storage.Adapter) storage.MDRetrieveQuery {
	return newRetrieve(conn(a))
}

// NewCreate opens a new createQuery query with the provided storage.Adapter.
func (e *Engine) NewCreate(a storage.Adapter) storage.MDCreateQuery {
	return newCreate(conn(a))
}

// NewDelete opens a new deleteQuery with the provided storage.Adapter.
func (e *Engine) NewDelete(a storage.Adapter) storage.MDDeleteQuery {
	return newDelete(conn(a))
}

// NewUpdate opens a new updateQuery with the provided storage.Adapter.
func (e *Engine) NewUpdate(a storage.Adapter) storage.MDUpdateQuery {
	return newUpdate(conn(a))
}

// NewMigrate opens a new migrateQuery with the provided storage.Adapter;
func (e *Engine) NewMigrate(a storage.Adapter) storage.MDMigrateQuery {
	return newMigrate(conn(a), e.cfg.Driver)
}
