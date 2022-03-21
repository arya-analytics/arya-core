package query

type Migrate struct {
	Base
}

func NewMigrate() *Migrate {
	c := &Migrate{}
	c.Base.Init(c)
	c.Model(&struct{}{})
	return c
}

func (m *Migrate) Verify() *Migrate {
	NewVerifyOpt(m.Pack())
	return m
}

func (m *Migrate) BindExec(e Execute) *Migrate {
	m.Base.BindExec(e)
	return m
}

// || VERIFY OPT ||

func NewVerifyOpt(p *Pack) {
	p.opts[verifyOptKey] = true
}

func VerifyOpt(p *Pack) bool {
	qo, ok := p.opts[verifyOptKey]
	if !ok {
		return false
	}
	return qo.(bool)
}
