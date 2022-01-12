package storage

import (
	"github.com/arya-analytics/aryacore/pkg/storage/engine"
)

// NewPooler creates a new Pooler.
func NewPooler(engines []engine.Base) *Pooler {
	return &Pooler{
		adapters: map[engine.Adapter]bool{},
		engines:  engines,
	}
}

// || POOLER ||

type Pooler struct {
	engines  []engine.Base
	adapters map[engine.Adapter]bool
}

// Retrieve retrieves an engine.Adapter based on the EngineType specified.
func (p *Pooler) Retrieve(r engine.Role) (a engine.Adapter, err error) {
	a, ok := p.findAdapter(r)
	if !ok {
		var err error
		a, err = p.newAdapter(r)
		if err != nil {
			return a, err
		}
		p.addAdapter(a)
	}
	return a, nil
}

func (p *Pooler) findAdapter(r engine.Role) (engine.Adapter, bool) {
	for a := range p.adapters {
		if a.Role() == r {
			return a, true
		}
	}
	return nil, false
}

func (p *Pooler) findEngine(r engine.Role) engine.Base {
	for _, e := range p.engines {
		if e.Role() == r {
			return e
		}
	}
	return nil
}

func (p *Pooler) newAdapter(r engine.Role) (engine.Adapter, error) {
	e := p.findEngine(r)
	a := e.NewAdapter()
	return a, nil
}

func (p *Pooler) addAdapter(a engine.Adapter) {
	p.adapters[a] = true
}
