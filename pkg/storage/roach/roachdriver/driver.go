package roachdriver

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

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
	// Port to Connect to at Host.
	Port int
	// Database to Connect to.
	Database string
	// Whether to open a TLS connection or not.
	UseTLS bool
}

func (d DriverPG) Connect() (*bun.DB, error) {
	c := d.buildConnector()
	db := sql.OpenDB(c)
	bunDB := bun.NewDB(db, pgdialect.New())
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
