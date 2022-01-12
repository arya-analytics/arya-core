package roach

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

// || DRIVER ||

type Driver int

const (
	DriverPG Driver = iota
	DriverSQLite
)

// || ENGINE ||

// Engine opens connections and execute queries with a roach database.
// implements the storage.MetaDataEngine interface.
type Engine struct {
	// Username for the database. Does not need to be specified if using DriverSQLite.
	Username string
	// Password for the database. Does not need to be specified if using DriverSQLite.
	Password string
	// Host IP for the database. Does not need to be specified if using DriverSQLite.
	Host     string
	// Port to connect to at Host. Does not need to be specified if using DriverSQLite.
	Port     int
	// Database to connect to. Does not need to be specified if using DriverSQLite.
	Database string
	// Whether to open a TLS connection or not.
	// Does not need to be specified if using DriverSQLite.
	UseTLS   bool
	// Driver is the connection driver used for the roach data store.
	// Current options are:
	// DriverPG which connects via the Postgres wire protocol.
	// DriverSQLite which uses an in memory SQLite database
	Driver   Driver
}

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapter() storage.Adapter {
	a := &adapter{
		id: uuid.New(),
		e:  e,
	}
	a.open()
	return a
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := e.bindAdapter(a)
	return ok
}

func (e *Engine) NewRetrieve(a storage.Adapter) storage.MetaDataRetrieve {
	ra, _ := e.bindAdapter(a)
	r := newRetrieve(e.conn(ra))
	return r
}

func (e *Engine) bindAdapter(a storage.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func (e *Engine) conn(a *adapter) *bun.DB {
	c, ok := a.conn().(*bun.DB)
	if !ok {
		log.Fatalln("Incorrect type specified")
	}
	return c
}

func (e *Engine) addr() string {
	return fmt.Sprintf("%s:%v", e.Host, e.Port)
}

func (e *Engine) tlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: e.UseTLS,
	}
}

// || ADAPTER ||

type adapter struct {
	id uuid.UUID
	db *bun.DB
	e  *Engine
}

// ID implements the storage.Adapter interface.
func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) conn() interface{} {
	return a.db
}

func (a *adapter) close() error {
	return a.db.Close()
}

func (a *adapter) open() {
	switch a.e.Driver {
	case DriverPG:
		a.db = pgConnect(a.e)
	case DriverSQLite:
		a.db = sqlLiteConnect()

	}
}

// || CONNECTORS ||

func pgConnect(e *Engine) *bun.DB {
	db := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithAddr(e.addr()),
			pgdriver.WithUser(e.Username),
			pgdriver.WithPassword(e.Password),
			pgdriver.WithDatabase(e.Database),
			pgdriver.WithTLSConfig(e.tlsConfig()),
		),
	)
	return bun.NewDB(db, pgdialect.New())
}

func sqlLiteConnect() *bun.DB {
	db, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		log.Fatalln(err)
	}
	return bun.NewDB(db, sqlitedialect.New())
}
