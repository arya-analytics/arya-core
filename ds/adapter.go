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
	configKey string
	conn      Conn
	getStatus GetStatus
	close     Close
}

func NewConnAdapter(configKey string, params ConnParams, connector Connector) ConnAdapter {
	conn, getStatus, _close := connector(params)
	return ConnAdapter{
		Id:        uuid.New(),
		created:   time.Now(),
		engine:    params.Engine,
		configKey: configKey,
		conn:      conn,
		getStatus: getStatus,
		close:     _close,
	}
}

func (ca *ConnAdapter) getExpired() bool {
	connExpirationTime := ca.created.Add(MaxConnAge)
	return connExpirationTime.Before(time.Now())
}
