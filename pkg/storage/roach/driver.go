package roach

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DriverRoach struct {
	// DSN is a connection string for the database. If specified,
	// all other fields except for Driver can be left blank.
	DSN string
	// Username for the database.
	Username string
	// Password for the database.
	Password string
	// Host IP for the database.
	Host string
	// Port to Connect to at Host.
	Port int
	// Database to Connect to.
	Database string
	// Whether to open a TLS connection or not.
	UseTLS bool
	// TransactionLogLevel
	TransactionLogLevel TransactionLogLevel
}

func (d DriverRoach) Connect() (*bun.DB, error) {
	c := d.buildConnector()
	db := sql.OpenDB(c)
	bunDB := bun.NewDB(db, pgdialect.New())
	setLogLevel(d.TransactionLogLevel, bunDB)
	return bunDB, nil
}

func (d DriverRoach) addr() string {
	return fmt.Sprintf("%s:%v", d.Host, d.Port)
}

func (d DriverRoach) buildConnector() *pgdriver.Connector {
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

type TransactionLogLevel int

const (
	// TransactionLogLevelNone logs no queries.
	TransactionLogLevelNone TransactionLogLevel = iota
	// TransactionLogLevelErr logs failed queries.
	TransactionLogLevelErr
	// TransactionLogLevelAll logs all queries.
	TransactionLogLevelAll
)

func setLogLevel(t TransactionLogLevel, db *bun.DB) {
	switch t {
	case TransactionLogLevelAll:
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	case TransactionLogLevelErr:
		db.AddQueryHook(bundebug.NewQueryHook())
	}
}
