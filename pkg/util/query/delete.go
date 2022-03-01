package query

// Delete deletes a model or set of models depending on the parameters passed.
// Delete requires that WherePK or WherePKs is called, and will panic upon execution if either isn't.
type Delete struct {
	where
}

func NewDelete() *Delete {
	d := &Delete{}
	d.baseInit(d)
	return d
}

// Model sets the model to be deleted. model must be passed as a pointer.
// Delete only uses the model for location lookups,
// and doesn't actually do anything  with the value.
// You're good to provide a nil pointer as long as it has the right type.
func (d *Delete) Model(m interface{}) *Delete {
	d.baseModel(m)
	return d
}

// WherePK queries by the primary of the model to be deleted.
func (d *Delete) WherePK(pk interface{}) *Delete {
	newPKOpt(d.Pack(), pk)
	return d
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (d *Delete) WherePKs(pks interface{}) *Delete {
	newPKsOpt(d.Pack(), pks)
	return d
}

func (d *Delete) BindExec(e Execute) *Delete {
	d.baseBindExec(e)
	return d
}
