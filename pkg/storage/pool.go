package storage

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"sync"
)

func NewPool() *Pool {
	return &Pool{
		adapters: map[internal.Adapter]bool{},
	}
}

type Pool struct {
	mu       sync.RWMutex
	adapters map[internal.Adapter]bool
}

func (p *Pool) Acquire(e internal.Engine) (a internal.Adapter, err error) {
	a, ok := p.findAdapter(e)
	if !ok {
		a, err = e.NewAdapter()
		p.addAdapter(a)
	}
	a.Acquire()
	return a, err
}

func (p *Pool) Release(a internal.Adapter) {
	a.Release()
}

func (p *Pool) findAdapter(e internal.Engine) (internal.Adapter, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for a := range p.adapters {
		if e.IsAdapter(a) && a.Healthy() {
			return a, true
		}
	}
	return nil, false
}

func (p *Pool) addAdapter(a internal.Adapter) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.adapters[a] = true
}
