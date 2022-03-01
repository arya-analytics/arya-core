package query

type Retrieve struct {
	where
}

// || CONSTRUCTOR ||

func NewRetrieve() *Retrieve {
	r := &Retrieve{}
	r.baseInit(r)
	return r
}

// || MODEL ||

func (r *Retrieve) Model(m interface{}) *Retrieve {
	r.baseModel(m)
	return r
}

// || WHERE ||

func (r *Retrieve) WherePK(pk interface{}) *Retrieve {
	r.wherePK(pk)
	return r
}

func (r *Retrieve) WherePKs(pks interface{}) *Retrieve {
	r.wherePKs(pks)
	return r
}

func (r *Retrieve) WhereFields(flds WhereFields) *Retrieve {
	r.whereFields(flds)
	return r
}

// || CALC ||

func (r *Retrieve) Calc(op Calc, fldName string, into interface{}) *Retrieve {
	newCalcOpt(r.Pack(), op, fldName, into)
	return r
}

// || FIELDS ||

func (r *Retrieve) Fields(flds ...string) *Retrieve {
	newFieldsOpt(r.Pack(), flds...)
	return r
}

func (r *Retrieve) Relation(rel string, flds ...string) *Retrieve {
	newRelationOpt(r.Pack(), rel, flds...)
	return r
}

// || EXEC ||

func (r *Retrieve) BindExec(e Execute) *Retrieve {
	r.baseBindExec(e)
	return r
}

// |||| OPTS ||||

// || RELATION ||

type RelationOpt struct {
	Rel    string
	Fields FieldsOpt
}

func newRelationOpt(p *Pack, rel string, flds ...string) {
	o := RelationOpt{rel, flds}
	_, ok := p.opts[relationOptKey]
	if !ok {
		p.opts[relationOptKey] = []RelationOpt{o}
	} else {
		p.opts[relationOptKey] = append(p.opts[relationOptKey].([]RelationOpt), o)
	}
}

func RelationOpts(p *Pack) []RelationOpt {
	o, ok := p.opts[relationOptKey]
	if !ok {
		return []RelationOpt{}
	}
	return o.([]RelationOpt)
}
