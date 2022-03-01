package query

// Update updates a model.
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

// WherePK queries the primary key of the model to be deleted.
func (u *Update) WherePK(pk interface{}) *Update {
	u.wherePK(pk)
	return u
}

// || FIELDS ||

// Fields marks the fields that need to be updated. If this option isn't specified,
// will replace all fields.
//
// NOTE: When calling Bulk, order matters. Fields must be called before Bulk.
func (u *Update) Fields(flds ...string) *Update {
	newFieldsOpt(u.Pack(), flds...)
	return u
}

// || BULK ||

// Bulk marks the update as a bulk update and allows for
// the update of multiple records. When Bulk updating,
// the primary key field of each model must be defined.
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
