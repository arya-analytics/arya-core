package storage

func newPooler() *pooler {
	return &pooler{
		adapters: map[Adapter]bool{},
	}
}

type pooler struct {
	adapters map[Adapter]bool
}

func (p *pooler) retrieve(e Engine) Adapter {
	a, ok := p.findAdapter(e)
	if !ok {
		a = p.newAdapter(e)
		p.addAdapter(a)
	}
	return a
}

func (p *pooler) findAdapter(e Engine) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *pooler) newAdapter(e Engine) Adapter {
	return e.NewAdapter()
}

func (p *pooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}
