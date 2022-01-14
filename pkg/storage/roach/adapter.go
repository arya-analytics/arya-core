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
)

// || ADAPTER ||

type adapter struct {
	id  uuid.UUID
	db  *bun.DB
	cfg Config
}

func newAdapter(cfg Config) *adapter {
	a := &adapter{
		id:  uuid.New(),
		cfg: cfg,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func conn(a storage.Adapter) *bun.DB {
	ra, ok := bindAdapter(a)
	if !ok {
		log.Fatalln("Couldn't bind roach adapter.")
	}
	return ra.conn()
}

// ID implements the storage.Adapter interface.
func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) conn() *bun.DB {
	return a.db
}

func (a *adapter) close() error {
	return a.db.Close()
}

func (a *adapter) open() {
	switch a.cfg.Driver {
	case DriverPG:
		a.db = pgConnect(a.cfg)
	case DriverSQLite:
		a.db = sqlLiteConnect()
	}
}

// || CONNECTORS ||

func pgConnect(cfg Config) *bun.DB {
	db := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithAddr(cfg.addr()),
			pgdriver.WithUser(cfg.Username),
			pgdriver.WithPassword(cfg.Password),
			pgdriver.WithDatabase(cfg.Database),
			pgdriver.WithTLSConfig(cfg.tls()),
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
