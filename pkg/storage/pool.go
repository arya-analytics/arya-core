package storage

func NewPool() *Pool {
	return &Pool{
		adapters: map[Adapter]bool{},
	}
}

type Pool struct {
	adapters map[Adapter]bool
}

func (p *Pool) Retrieve(e Engine) Adapter {
	a, ok := p.findAdapter(e)
	if !ok {
		a = p.newAdapter(e)
		p.addAdapter(a)
	}
	return a
}

func (p *Pool) findAdapter(e Engine) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *Pool) newAdapter(e Engine) Adapter {
	return e.NewAdapter()
}

func (p *Pool) addAdapter(a Adapter) {
	p.adapters[a] = true
}
