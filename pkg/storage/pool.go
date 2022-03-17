package storage

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"sync"
)

type adapterState struct {
	Demand int
}

func (as *adapterState) Acquire() {
	as.Demand += 1
}

func (as *adapterState) Release() {
	as.Demand -= 1
}

func NewPool() *Pool {
	return &Pool{
		adapters: map[internal.Adapter]*adapterState{},
	}
}

type Pool struct {
	mu       sync.RWMutex
	adapters map[internal.Adapter]*adapterState
}

func (p *Pool) Retrieve(e internal.Engine) (a internal.Adapter, err error) {
	a, ok := p.findAdapter(e)
	if !ok {
		a, err = e.NewAdapter()
		p.addAdapter(a)
	}
	return a, err
}

func (p *Pool) Release(a internal.Adapter) {
	p.adapters[a].Release()
}

func (p *Pool) findAdapter(e internal.Engine) (internal.Adapter, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for a, s := range p.adapters {
		if e.IsAdapter(a) && a.DemandCap() > s.Demand {
			s.Acquire()
			return a, true
		}
	}
	return nil, false
}

func (p *Pool) addAdapter(a internal.Adapter) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.adapters[a] = &adapterState{}
}
