package roach

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Driver interface {
	Connect() (*bun.DB, error)
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
