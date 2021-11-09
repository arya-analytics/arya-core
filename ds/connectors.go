package ds

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Connector func(params Config) (conn Conn)

func getConnector(engine string) Connector {
	connectors := map[string]Connector{
		"github.com/uptrace/bun/driver/pgdriver": pgConnector,
	}
	fmt.Println(engine)
	return connectors[engine]
}

// || POSTGRES ||
func pgConnector(params Config) (conn Conn) {
	dsn := "postgres://" + params.User + ":" + params.Password + "@" + params.Host +
		"/" + params.Name
	if !params.Secure {
		dsn += "?sslmode=disable"
	}
	baseDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(baseDB, pgdialect.New())
}
