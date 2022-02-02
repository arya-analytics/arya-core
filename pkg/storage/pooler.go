package storage

func newPooler() *pooler {
	return &pooler{
		adapters: map[Adapter]bool{},
	}
}

// || POOLER ||

type pooler struct {
	adapters map[Adapter]bool
}

// Retrieve retrieves an engine.Adapter based on the EngineType specified.
func (p *pooler) retrieve(e BaseEngine) Adapter {
	a, ok := p.findAdapter(e)
	if !ok {
		a = p.newAdapter(e)
		p.addAdapter(a)
	}
	return a
}

func (p *pooler) findAdapter(e BaseEngine) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *pooler) newAdapter(e BaseEngine) Adapter {
	return e.NewAdapter()
}

func (p *pooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}
