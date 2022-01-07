package ds

import (
	"github.com/google/uuid"
	"time"
)

type Conn interface{}

type ConnAdapter struct {
	Id      uuid.UUID
	created time.Time
	engine  Engine
	key     string
	conn    Conn
}

func NewConnAdapter(key string, params Config, connector Connector) ConnAdapter {
	conn := connector(params)
	return ConnAdapter{
		Id:      uuid.New(),
		created: time.Now(),
		engine:  params.Engine,
		key:     key,
		conn:    conn,
	}
}
