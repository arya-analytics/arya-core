package storage

// NewPooler creates a new Pooler.
func NewPooler() *Pooler {
	return &Pooler{
		adapters: map[Adapter]bool{},
	}
}

// || POOLER ||

type Pooler struct {
	adapters map[Adapter]bool
}

// Retrieve retrieves an engine.Adapter based on the EngineType specified.
func (p *Pooler) Retrieve(e EngineBase) (a Adapter, err error) {
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

func (p *Pooler) findAdapter(e EngineBase) (Adapter, bool) {
	for a := range p.adapters {
		if e.IsAdapter(a) {
			return a, true
		}
	}
	return nil, false
}

func (p *Pooler) newAdapter(e EngineBase) (Adapter, error) {
	a := e.NewAdapter()
	return a, nil
}

func (p *Pooler) addAdapter(a Adapter) {
	p.adapters[a] = true
}
