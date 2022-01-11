package storage

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Adapter interface {
	ID() uuid.UUID
	Release() error
	Status() ConnStatus
	Conn() interface{}
	Role() EngineRole
	close() error
	open() error
}

type ConnStatus int

const (
	ConnStatusReady ConnStatus = iota
)

func NewPooler(cfgChain ConfigChain) *Pooler {
	return &Pooler{cfgChain: cfgChain, adapters: []Adapter{}}
}

type Pooler struct {
	cfgChain ConfigChain
	adapters []Adapter
}

func (p *Pooler) Retrieve(r EngineRole) (a Adapter) {
	a, ok := p.findAdapter(r)
	if !ok {
		var err error
		a, err = p.newAdapter(r)
		if err != nil {
			log.Fatalln(err)
		}
		p.addAdapter(a)
	}
	return a
}

func (p *Pooler) findAdapter(r EngineRole) (Adapter, bool) {
	for _, a := range p.adapters {
		if a.Role() == r {
			return a, true
		}
	}

	return nil, false
}

func (p *Pooler) newAdapter(r EngineRole) (Adapter, error) {
	cfg, err := p.cfgChain.Retrieve(r)
	if err != nil {
		log.Fatalln(err)
	}
	a, err := NewAdapter(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return a, nil
}

func (p *Pooler) addAdapter(a Adapter) {
	p.adapters = append(p.adapters, a)
}

func NewAdapter(cfg Config) (Adapter, error) {
	switch cfg.Type() {
	case EngineTypeMDStub:
		return NewMDStubAdapter(cfg)
	}
	return nil, nil
}

type MDStubConn struct {}

type MDStubAdapter struct {
	id uuid.UUID
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

func (a *MDStubAdapter) Status() ConnStatus {
	return ConnStatusReady
}

func (a *MDStubAdapter) Role() EngineRole {
	return EngineRoleMetaData
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
