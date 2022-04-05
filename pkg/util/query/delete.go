package query

// Delete deletes a model or set of models depending on the parameters passed.
// Delete requires that WherePK or WherePKs is called, and will panic upon execution if either isn't.
type Delete struct {
	Where
}

// NewDelete opens a new Delete query.
func NewDelete() *Delete {
	d := &Delete{}
	d.Base.Init(d)
	return d
}

// Model sets the model to be deleted. model must be passed as a pointer.
// Delete only uses the model for location lookups,
// and doesn't actually do anything  with the value.
// You're good to provide a nil pointer as long as it has the right type.
func (d *Delete) Model(m interface{}) *Delete {
	d.Base.Model(m)
	return d
}

// WherePK queries by the primary of the model to be deleted.
func (d *Delete) WherePK(pk interface{}) *Delete {
	d.Where.WherePK(pk)
	return d
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (d *Delete) WherePKs(pks interface{}) *Delete {
	d.Where.WherePKs(pks)
	return d
}

// WhereFields queries by a set of key value pairs where the key represents a field name
// and the value represents a value to match with.
func (d *Delete) WhereFields(flds WhereFields) *Delete {
	d.Where.WhereFields(flds)
	return d
}

// BindExec binds Execute that Delete will use to run the query.
// This method MUST be called before calling Exec.
func (d *Delete) BindExec(e Execute) *Delete {
	d.Base.BindExec(e)
	return d
}
