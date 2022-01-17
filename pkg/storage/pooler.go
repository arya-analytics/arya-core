package storage

func UnsafeNewPooler() *pooler {
	return newPooler()
}

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
func (p *pooler) Retrieve(e BaseEngine) (a Adapter, err error) {
	a, ok := p.findAdapter(e)
	if !ok {
		var err error
		a, err = p.newAdapter(e)
		if err != nil {
			return a, err
		}
		p.addAdapter(a)
	}
	return a, nil
}

func (p *pooler) findAdapter(e BaseEngine) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *pooler) newAdapter(e BaseEngine) (Adapter, error) {
	a := e.NewAdapter()
	return a, nil
}

func (p *pooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}
