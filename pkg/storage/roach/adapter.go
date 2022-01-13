package roach

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

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
