package storage

func newPool() *pool {
	return &pool{
		adapters: map[Adapter]bool{},
	}
}

type pool struct {
	adapters map[Adapter]bool
}

func (p *pool) retrieve(e Engine) Adapter {
	a, ok := p.findAdapter(e)
	if !ok {
		a = p.newAdapter(e)
		p.addAdapter(a)
	}
	return a
}

func (p *pool) findAdapter(e Engine) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *pool) newAdapter(e Engine) Adapter {
	return e.NewAdapter()
}

func (p *pool) addAdapter(a Adapter) {
	p.adapters[a] = true
}
