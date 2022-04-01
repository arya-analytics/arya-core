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
	NewAdapt() (Adapt[T], error)
	Match(T) bool
}

type Pool[T any] struct {
	mu        sync.RWMutex
	Factories map[AdaptFactory[T]]bool
	Adapters  map[Adapt[T]]bool
}

func New[T any]() *Pool[T] {
	return &Pool[T]{
		Factories: make(map[AdaptFactory[T]]bool),
		Adapters:  make(map[Adapt[T]]bool),
	}
}

func (p *Pool[T]) Acquire(match T) (a Adapt[T], err error) {
	a, ok := p.findAdapter(match)
	if !ok {
		a, err = p.findFactory(match).NewAdapt()
		p.addAdapt(a)
	}
	a.Acquire()
	return a, err
}

func (p *Pool[T]) AddFactory(f AdaptFactory[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Factories[f] = true
}

func (p *Pool[T]) Release(a Adapt[T]) {
	a.Release()
}

func (p *Pool[T]) findFactory(match T) AdaptFactory[T] {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for f := range p.Factories {
		if f.Match(match) {
			return f
		}
	}
	panic("no factory could be found for pool")
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
