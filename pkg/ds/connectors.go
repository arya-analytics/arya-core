package ds

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Connector func(params Config) (conn Conn)

func getConnector(engine Engine) Connector {
	connectors := map[Engine]Connector{
		Postgres:  pgConnector,
		GorillaWS: gorillaWSConnector,
	}
	fmt.Println(engine)
	return connectors[engine]
}

// || POSTGRES ||
func pgConnector(cfg Config) (conn Conn) {
	dsn := "postgres://" + cfg.Auth.User + ":" + cfg.Auth.Password + "@" + cfg.
		Host +
		"/" + cfg.Name
	if !cfg.Secure {
		dsn += "?sslmode=disable"
	}
	baseDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(baseDB, pgdialect.New())
}

// || WEBSOCKET ||
func gorillaWSConnector(cfg Config) (conn Conn) {
	dsn := "ws://" + cfg.Host + ":" + cfg.Port + cfg.Name
	conn, _, err := websocket.DefaultDialer.Dial(dsn, nil)
	if err != nil {
		panic(err)
	}
	return conn
}
