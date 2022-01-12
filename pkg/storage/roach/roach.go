package roach

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/engine"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const (
	role = engine.RoleMetaData
)

// || DRIVER ||

type Driver int

const (
	DriverPG Driver = iota
	DriverSQLite
)

// || ENGINE ||

type Engine struct {
	Username string
	Password string
	Host     string
	Database string
	Port     int
	UseTLS   bool
	Driver   Driver
}

func (e *Engine) Role() engine.Role {
	return role
}

func (e *Engine) NewAdapter() (a *adapter, err error) {
	a = &adapter{
		id: uuid.New(),
		e:  e,
	}
	a.open()
	return a, err
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

func (a *adapter) Role() engine.Role {
	return a.e.Role()
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
