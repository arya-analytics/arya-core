package roach

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
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

// Migrate migrates the database.
func (e *Engine) Migrate(ctx context.Context, a storage.Adapter) error {
	ra, _ := e.bindAdapter(a)
	db := e.conn(ra)
	m := newMigrator(db)
	err := m.init(ctx)
	if err != nil {
		return err
	}
	err = m.migrate(ctx)
	return err
}

// VerifyMigrations verifies that the migrations were executed correctly
func (e *Engine) VerifyMigrations(ctx context.Context, a storage.Adapter) error {
	ra, _ := e.bindAdapter(a)
	db := e.conn(ra)
	m := newMigrator(db)
	return m.verify(ctx)
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



