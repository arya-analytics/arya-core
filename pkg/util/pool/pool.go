package pool

import (
	"sync"
)

type Adapt[T any] interface {
	Healthy() bool
	Acquire()
	Release()
	Match(T) bool
}

type AdaptFactory[T any] interface {
	NewAdapt(T) (Adapt[T], error)
}

type Pool[T any] struct {
	mu       sync.RWMutex
	Factory  AdaptFactory[T]
	Adapters map[Adapt[T]]bool
}

func New[T any]() *Pool[T] {
	return &Pool[T]{
		Adapters: make(map[Adapt[T]]bool),
	}
}

func (p *Pool[T]) Acquire(match T) (a Adapt[T], err error) {
	a, ok := p.findAdapter(match)
	if !ok {
		a, err = p.Factory.NewAdapt(match)
		p.addAdapt(a)
	}
	a.Acquire()
	return a, err
}

func (p *Pool[T]) Release(a Adapt[T]) {
	a.Release()
}

func (p *Pool[T]) findAdapter(match T) (Adapt[T], bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for a := range p.Adapters {
		if a.Match(match) && a.Healthy() {
			return a, true
		}
	}
	return nil, false
}

func (p *Pool[T]) addAdapt(a Adapt[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Adapters[a] = true
}
