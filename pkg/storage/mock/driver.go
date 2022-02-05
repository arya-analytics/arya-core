package mock

import (
	"database/sql"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// |||| DRIVER PG ||||

type DriverPG struct{}

func NewDriverPG() DriverPG {
	return DriverPG{}
}

func (d DriverPG) Connect() (*bun.DB, error) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		return nil, err
	}
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(ts.PGURL().String())))
	db := bun.NewDB(sqlDB, pgdialect.New())
	return db, nil
}
