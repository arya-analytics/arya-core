package ds

import (
	"github.com/google/uuid"
	"time"
)

type ConnStatus int

const (
	Active ConnStatus = iota
	Inactive
)

type Conn interface{}
type GetStatus func(conn Conn) ConnStatus
type Close func(conn Conn)

type ConnAdapter struct {
	Id        uuid.UUID
	created   time.Time
	engine    string
	key       string
	conn      Conn
	getStatus GetStatus
	close     Close
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