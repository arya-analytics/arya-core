package roach

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Driver interface {
	Connect() (*bun.DB, error)
}

type DriverPG struct {
	// DSN is a connection string for the database. If specified,
	// all other fields except for Driver can be left blank.
	DSN string
	// Username for the database.
	Username string
	// Password for the database.
	Password string
	// Host IP for the database.
	Host string
	// Port to connect to at Host.
	Port int
	// Database to connect to.
	Database string
	// Whether to open a TLS connection or not.
	UseTLS bool
	// TransactionLogLevel is the log level for executed SQL queries
	TransactionLogLevel TransactionLogLevel
}

func NewDriverPG() DriverPG {
	return DriverPG{}
}

func (d DriverPG) Connect() (*bun.DB, error) {
	c := d.buildConnector()
	db := sql.OpenDB(c)
	bunDB := bun.NewDB(db, pgdialect.New())
	setLogLevel(d.TransactionLogLevel, bunDB)
	return bunDB, nil
}

func (d DriverPG) addr() string {
	return fmt.Sprintf("%s:%v", d.Host, d.Port)
}

func (d DriverPG) buildConnector() *pgdriver.Connector {
	if d.DSN != "" {
		return pgdriver.NewConnector(pgdriver.WithDSN(d.DSN))
	}
	return pgdriver.NewConnector(
		pgdriver.WithAddr(d.addr()),
		pgdriver.WithInsecure(d.UseTLS),
		pgdriver.WithUser(d.Username),
		pgdriver.WithPassword(d.Password),
		pgdriver.WithDatabase(d.Database))
}

type DriverSQLite struct {
	// TransactionLogLevel is the log level for executed SQL queries
	TransactionLogLevel TransactionLogLevel
}

func (d DriverSQLite) Connect() (*bun.DB, error) {
	db, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	bunDB := bun.NewDB(db, sqlitedialect.New())
	setLogLevel(d.TransactionLogLevel, bunDB)
	return bunDB, nil
}

func setLogLevel(t TransactionLogLevel, db *bun.DB) {
	switch t {
	case TransactionLogLevelAll:
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	case TransactionLogLevelErr:
		db.AddQueryHook(bundebug.NewQueryHook())
	}
}
