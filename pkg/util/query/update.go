package query

// Update updates a model.
type Update struct {
	Where
}

// || CONSTRUCTOR ||

// NewUpdate opens a new Update query.
func NewUpdate() *Update {
	u := &Update{}
	u.Base.Init(u)
	return u
}

// || MODEL ||

// Model sets the model to bind the results into. Model must be passed as a pointer or a *model.Reflect.
func (u *Update) Model(m interface{}) *Update {
	u.Base.Model(m)
	return u
}

// || WHERE ||

// WherePK queries the primary key of the model to be deleted.
func (u *Update) WherePK(pk interface{}) *Update {
	u.Where.WherePK(pk)
	return u
}

// || FIELDS ||

// Fields marks the fields that need to be updated. If this option isn't specified,
// will replace all fields.
//
// NOTE: When calling Bulk, order matters. Fields must be called before Bulk.
func (u *Update) Fields(fields ...string) *Update {
	NewFieldsOpt(u.Pack(), fields...)
	return u
}

// || BULK ||

// Bulk marks the update as a bulk update and allows for
// the update of multiple records. When Bulk updating,
// the primary key field of each model must be defined.
func (u *Update) Bulk() *Update {
	NewBulkUpdateOpt(u.Pack())
	return u
}

// || EXEC ||

// BindExec binds Execute that Update will use to run the query.
// This method MUST be called before calling Exec.
func (u *Update) BindExec(e Execute) *Update {
	u.Base.BindExec(e)
	return u
}

// |||| OPTS |||

// || BULK ||

// NewBulkUpdateOpt creates a new BulkUpdateOpt.
func NewBulkUpdateOpt(p *Pack) {
	p.opts[bulkUpdateOptKey] = true
}

// BulkUpdateOpt returns true if the Update is a bulk update.
func BulkUpdateOpt(p *Pack) bool {
	_, ok := p.Query().(*Update)
	if !ok {
		panic("can't retrieve a bulk query opt from non bulk query")
	}
	bulkOpt, ok := p.RetrieveOpt(bulkUpdateOptKey)
	if !ok {
		return false
	}
	return bulkOpt.(bool)
}
