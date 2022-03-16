package query

// FieldsOpt holds a slice of strings representing the model fields to retrieve or execute on in a Query.
type FieldsOpt []string

// NewFieldsOpt creates a new FieldsOpt and binds it to the provided Pack.
func NewFieldsOpt(p *Pack, flds ...string) {
	p.opts[fieldsOptKey] = FieldsOpt(flds)
}

// RetrieveFieldsOpt retrieves the FieldsOpt from Pack p. Returns false if Pack does not have a FieldsOpt specified.
func RetrieveFieldsOpt(p *Pack) (FieldsOpt, bool) {
	qo, ok := p.opts[fieldsOptKey]
	if !ok {
		return FieldsOpt{}, false
	}
	return qo.(FieldsOpt), true
}

// ContainsAny returns true if FieldsOpt contains any of the provided fields.
func (fo FieldsOpt) ContainsAny(flds ...string) (contains bool) {
	for _, qFld := range flds {
		for _, fld := range fo {
			if qFld == fld {
				contains = true
			}
		}
	}
	return contains
}

// ContainsAll returns true if FieldsOpt contains all provided fields.
func (fo FieldsOpt) ContainsAll(flds ...string) bool {
	for _, fld := range flds {
		present := false
		for _, fqoFld := range fo {
			if fld == fqoFld {
				present = true
			}
		}
		if !present {
			return false
		}
	}
	return true
}

// AllExcept returns a new FieldsOpt with all the same fields except for the provided flds.
func (fo FieldsOpt) AllExcept(flds ...string) (filteredFqo FieldsOpt) {
	for _, fld := range fo {
		for _, eFld := range flds {
			if eFld != fld {
				filteredFqo = append(filteredFqo, fld)
			}
		}
	}
	return filteredFqo
}

// Append returns a new FieldsOpt with the provided flds appended to it.
// NOTE: Will remove duplicates.
func (fo FieldsOpt) Append(flds ...string) (nFo FieldsOpt) {
	nFo = fo
	for _, fld := range flds {
		if !nFo.ContainsAll(fld) {
			nFo = append(nFo, fld)
		}
	}
	return nFo
}
