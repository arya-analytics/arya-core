package storage

import "sync"

type AdapterState struct {
	Demand int
}

func (as *AdapterState) Acquire() {
	as.Demand += 1
}

func (as *AdapterState) Release() {
	as.Demand -= 1
}

func NewPool() *Pool {
	return &Pool{
		adapters: map[Adapter]*AdapterState{},
	}
}

type Pool struct {
	mu       sync.RWMutex
	adapters map[Adapter]*AdapterState
}

func (p *Pool) Retrieve(e Engine) (a Adapter, err error) {
	a, ok := p.findAdapter(e)
	if !ok {
		a, err = e.NewAdapter()
		p.addAdapter(a)
	}
	return a, err
}

func (p *Pool) Release(a Adapter) {
	p.adapters[a].Release()
}

func (p *Pool) findAdapter(e Engine) (Adapter, bool) {
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

func (p *Pool) addAdapter(a Adapter) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.adapters[a] = &AdapterState{}
}
