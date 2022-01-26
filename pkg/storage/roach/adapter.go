package roach

import (
	"database/sql"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// TODO: Fix TLS config
// TODO: Test connection issues

// || ADAPTER ||

type Adapter struct {
	id  uuid.UUID
	db  *bun.DB
	cfg Config
}

func newAdapter(cfg Config) *Adapter {
	a := &Adapter{
		id:  uuid.New(),
		cfg: cfg,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*Adapter, bool) {
	ra, ok := a.(*Adapter)
	return ra, ok
}

func conn(a storage.Adapter) *bun.DB {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind roach adapter.")
	}
	return ra.conn()
}

// ID implements the storage.Adapter interface.
func (a *Adapter) ID() uuid.UUID {
	return a.id
}

func (a *Adapter) conn() *bun.DB {
	return a.db
}

func (a *Adapter) setLogLevel() {
	switch a.cfg.TransactionLogLevel {
	case TransactionLogLevelAll:
		a.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	case TransactionLogLevelErr:
		a.db.AddQueryHook(bundebug.NewQueryHook())
	}
}

func (a *Adapter) open() {
	switch a.cfg.Driver {
	case DriverPG:
		a.db = connectToPG(a.cfg)
	case DriverSQLite:
		a.db = connectToSqlite()
	}
	a.setLogLevel()
}

// || CONNECTORS ||

func pgConfig(cfg Config) *pgdriver.Connector {
	if cfg.DSN != "" {
		return pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN))
	}
	return pgdriver.NewConnector(
		pgdriver.WithAddr(cfg.addr()),
		pgdriver.WithInsecure(cfg.UseTLS),
		pgdriver.WithUser(cfg.Username),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithDatabase(cfg.Database))
}

func connectToPG(cfg Config) *bun.DB {
	db := sql.OpenDB(pgConfig(cfg))
	return bun.NewDB(db, pgdialect.New())
}

func connectToSqlite() *bun.DB {
	db, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		log.Fatalln(err)
	}
	return bun.NewDB(db, sqlitedialect.New())
}
