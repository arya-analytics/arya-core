package query

// Retrieve retrieves a model or set of models depending on the parameters passed.
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

// Model sets the model to bind the results into. model must be passed as a pointer.
// If you're expecting multiple return values,
// pass a pointer to a slice. If you're expecting one return value,
// pass a struct. NOTE: If a struct is passed, and multiple values are returned,
// the struct is assigned to the value of the first result.
func (r *Retrieve) Model(m interface{}) *Retrieve {
	r.baseModel(m)
	return r
}

// || WHERE ||

// WherePK queries by the primary of the model to be deleted.
func (r *Retrieve) WherePK(pk interface{}) *Retrieve {
	r.wherePK(pk)
	return r
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (r *Retrieve) WherePKs(pks interface{}) *Retrieve {
	r.wherePKs(pks)
	return r
}

// WhereFields queries by a set of key value pairs where the key represents a field name
// and the value represent a value to match with.
func (r *Retrieve) WhereFields(flds WhereFields) *Retrieve {
	r.whereFields(flds)
	return r
}

// || CALC ||

// Calc executes a calculation on the specified field. It binds the calculation
// into the argument 'into'. See Calculate for available calculations.
func (r *Retrieve) Calc(op Calc, fldName string, into interface{}) *Retrieve {
	newCalcOpt(r.Pack(), op, fldName, into)
	return r
}

// || FIELDS ||

// Fields retrieves only the fields specified.
func (r *Retrieve) Fields(flds ...string) *Retrieve {
	newFieldsOpt(r.Pack(), flds...)
	return r
}

// Relation retrieves the specified fields from the relation.
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
