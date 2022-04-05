package query

// Migrate migrates the schema of a data store.
type Migrate struct {
	Base
}

// NewMigrate creates a new Migrate query.
func NewMigrate() *Migrate {
	c := &Migrate{}
	c.Base.Init(c)
	c.Model(&struct{}{})
	return c
}

// Verify verifies that the schema of the data store is up-to-date.
func (m *Migrate) Verify() *Migrate {
	NewVerifyOpt(m.Pack())
	return m
}

// BindExec binds Execute that Migrate will use to run the query.
// This method must be called before calling Exec.
func (m *Migrate) BindExec(e Execute) *Migrate {
	m.Base.BindExec(e)
	return m
}

// || VERIFY OPT ||

// NewVerifyOpt creates a new VerifyOpt.
func NewVerifyOpt(p *Pack) {
	p.SetOpt(verifyOptKey, true)
}

// VerifyOpt is an option to the Pack indicating that Migrate query
// should verify migrations are up-to-date instead of running the migrations themselves.
func VerifyOpt(p *Pack) bool {
	qo, ok := p.RetrieveOpt(verifyOptKey)
	if !ok {
		return false
	}
	return qo.(bool)
}
