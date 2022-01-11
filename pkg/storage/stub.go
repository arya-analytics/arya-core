package storage

import "github.com/google/uuid"

// || META DATA STUBS ||

type MDStubConn struct{}

type MDStubAdapter struct {
	id   uuid.UUID
	cfg  Config
	conn *MDStubConn
}

func NewMDStubAdapter(cfg Config) (a *MDStubAdapter, err error) {
	a = &MDStubAdapter{cfg: cfg, conn: &MDStubConn{}, id: uuid.New()}
	err = a.open()
	return a, err
}

func (a *MDStubAdapter) ID() uuid.UUID {
	return a.id
}

func (a *MDStubAdapter) Release() error {
	return nil
}

func (a *MDStubAdapter) Status() AdapterStatus {
	return ConnStatusReady
}

func (a *MDStubAdapter) Type() EngineType {
	return EngineTypeMDStub
}

func (a *MDStubAdapter) close() error {
	return nil
}

func (a *MDStubAdapter) open() error {
	return nil
}

func (a *MDStubAdapter) Conn() interface{} {
	return a.conn
}