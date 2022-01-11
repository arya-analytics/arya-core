package storage

import (
	"fmt"
	"github.com/google/uuid"
)

// || ADAPTER ||

type Adapter interface {
	ID() uuid.UUID
	Status() AdapterStatus
	Conn() interface{}
	Type() EngineType
	close() error
	open() error
}

type AdapterStatus int

const (
	ConnStatusReady AdapterStatus = iota
)

// || POOLER ERROR ||

type PoolerError struct {
	Op string
	Et EngineType
}

func (p PoolerError) Error() string {
	return fmt.Sprintf("%s %v", p.Op, p.Et)
}

func NewPoolerError(op string, Et EngineType) PoolerError {
	return PoolerError{Op: op, Et: Et}
}

func NewPooler(cfgChain ConfigChain) *Pooler {
	return &Pooler{cfgChain: cfgChain, adapters: map[Adapter]bool{}}
}

// || POOLER ||

type Pooler struct {
	cfgChain ConfigChain
	adapters map[Adapter]bool
}

func (p *Pooler) Retrieve(et EngineType) (a Adapter, err error) {
	a, ok := p.findAdapter(et)
	if !ok {
		var err error
		a, err = p.newAdapter(et)
		if err != nil {
			return a, err
		}
		p.addAdapter(a)
	}
	return a, nil
}

func (p *Pooler) findAdapter(et EngineType) (Adapter, bool) {
	for a := range p.adapters {
		if a.Type() == et {
			return a, true
		}
	}
	return nil, false
}

func (p *Pooler) newAdapter(et EngineType) (Adapter, error) {
	cfg, err := p.cfgChain.Retrieve(et)
	if err != nil {
		return nil, err
	}
	a, err := NewAdapter(cfg)
	if err != nil {
		return a, err
	}
	return a, nil
}

func (p *Pooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}

func NewAdapter(cfg Config) (Adapter, error) {
	et := cfg.Type()
	switch et {
	case EngineTypeMDStub:
		return NewMDStubAdapter(cfg)
	default:
		return nil, NewPoolerError("adapter type does not exist", et)
	}
}
