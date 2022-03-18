package internal

import (
	"sync"
)

func NewPool() *Pool {
	return &Pool{
		adapters: map[Adapter]bool{},
	}
}

type Pool struct {
	mu       sync.RWMutex
	adapters map[Adapter]bool
}

func (p *Pool) Acquire(e Engine) (a Adapter, err error) {
	a, ok := p.findAdapter(e)
	if !ok {
		a, err = e.NewAdapter()
		p.addAdapter(a)
	}
	a.Acquire()
	return a, err
}

func (p *Pool) Release(a Adapter) {
	a.Release()
}

func (p *Pool) findAdapter(e Engine) (Adapter, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for a := range p.adapters {
		if e.IsAdapter(a) && a.Healthy() {
			return a, true
		}
	}
	return nil, false
}

func (p *Pool) addAdapter(a Adapter) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.adapters[a] = true
}
