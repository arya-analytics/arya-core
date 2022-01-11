package storage

import "log"

type Pooler interface {
	Retrieve(et EngineType) Adapter
}

type Adapter interface {
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

type DefaultPooler struct {
	cfgChain ConfigChain
	adapters map[Adapter]bool
}

func (p *DefaultPooler) Retrieve(r EngineRole) (a Adapter) {
	a, ok := p.findAdapter(r)
	if !ok {
		a, err := p.newAdapter(r)
		if err != nil {
			log.Fatalln(err)
		}
		p.addAdapter(a)
	}
	return a
}

func (p *DefaultPooler) findAdapter(r EngineRole) (Adapter, bool) {
	for a := range p.adapters {
		if a.Role() == r {
			return a, true
		}
	}
	return nil, false
}

func (p *DefaultPooler) newAdapter(r EngineRole) (Adapter, error) {
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

func (p *DefaultPooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}

func NewAdapter(cfg Config) (Adapter, error) {
	switch cfg.Type() {
	case EngineTypeMDStub:
		return NewMDStubAdapter(cfg)
	}
	return nil, nil
}

type MDStubAdapter struct {
	cfg Config
}

func NewMDStubAdapter(cfg Config) (a MDStubAdapter, err error){
	a = MDStubAdapter{cfg}
	err = a.open()
	return a, err
}

func (a MDStubAdapter) Release() error {
	return nil
}

func (a MDStubAdapter) Status() ConnStatus {
	return ConnStatusReady
}

func (a MDStubAdapter) Role() EngineRole {
	return EngineRoleMetaData
}

func (a MDStubAdapter) close() error {
	return nil
}

func (a MDStubAdapter) open() error {
	return nil
}

func (a MDStubAdapter) Conn() interface{} {
	return nil
}
