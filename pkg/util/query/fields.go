package query

func RetrieveFieldsOpt(p *Pack) (FieldsOpt, bool) {
	qo, ok := p.opts[fieldsOptKey]
	if !ok {
		return FieldsOpt{}, false
	}
	return qo.(FieldsOpt), true
}

type FieldsOpt []string

func newFieldsOpt(p *Pack, flds ...string) {
	p.opts[fieldsOptKey] = FieldsOpt(flds)
}

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

func (fo FieldsOpt) Append(flds ...string) (nFo FieldsOpt) {
	nFo = fo
	for _, fld := range flds {
		if !nFo.ContainsAll(fld) {
			nFo = append(nFo, fld)
		}
	}
	return nFo
}
