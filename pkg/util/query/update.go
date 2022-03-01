package query

type Update struct {
	where
}

// || CONSTRUCTOR ||

func NewUpdate() *Update {
	u := &Update{}
	u.baseInit(u)
	return u
}

// || MODEL ||

func (u *Update) Model(m interface{}) *Update {
	u.baseModel(m)
	return u
}

// || WHERE ||

func (u *Update) WherePK(pk interface{}) *Update {
	u.wherePK(pk)
	return u
}

// || BULK ||

func (u *Update) Bulk() *Update {
	newBulkUpdateOpt(u.Pack())
	return u
}

// || EXEC ||

func (u *Update) BindExec(e Execute) *Update {
	u.baseBindExec(e)
	return u
}

// |||| OPTS |||

// || BULK ||

func newBulkUpdateOpt(p *Pack) {
	p.opts[bulkUpdateOptKey] = true
}

func BulkUpdateOpt(p *Pack) bool {
	_, ok := p.Query().(*Update)
	if !ok {
		panic("can't retrieve a bulk query opt from non bulk query")
	}
	bulkOpt, ok := p.opts[bulkUpdateOptKey]
	if !ok {
		return false
	}
	return bulkOpt.(bool)
}
