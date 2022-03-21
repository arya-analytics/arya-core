package query

// Retrieve retrieves a model or set of models depending on the parameters passed.
type Retrieve struct {
	Where
}

// || CONSTRUCTOR ||

func NewRetrieve() *Retrieve {
	r := &Retrieve{}
	r.Base.Init(r)
	return r
}

// || MODEL ||

// Model sets the model to bind the results into. model must be passed as a pointer.
// If you're expecting multiple return values,
// pass a pointer to a slice. If you're expecting one return value,
// pass a struct. NOTE: If a struct is passed, and multiple values are returned,
// the struct is assigned to the value of the first result.
func (r *Retrieve) Model(m interface{}) *Retrieve {
	r.Base.Model(m)
	return r
}

// || WHERE ||

// WherePK queries by the primary of the model to be deleted.
func (r *Retrieve) WherePK(pk interface{}) *Retrieve {
	r.Where.WherePK(pk)
	return r
}

// WherePKs queries by a set of primary keys of models to be deleted.
func (r *Retrieve) WherePKs(pks interface{}) *Retrieve {
	r.Where.WherePKs(pks)
	return r
}

// WhereFields queries by a set of key value pairs where the key represents a field name
// and the value represent a value to match with.
func (r *Retrieve) WhereFields(flds WhereFields) *Retrieve {
	r.Where.WhereFields(flds)
	return r
}

// || CALC ||

// Calc executes a calculation on the specified field. It binds the calculation
// into the argument 'into'. See Calculate for available calculations.
func (r *Retrieve) Calc(op Calc, fld string, into interface{}) *Retrieve {
	NewCalcOpt(r.Pack(), op, fld, into)
	return r
}

// || ORDER ||

// Order orders the results in a specific direction.
func (r *Retrieve) Order(order Order, fld string) *Retrieve {
	NewOrderOpt(r.Pack(), order, fld)
	return r
}

// Limit limits the number of results returned.
func (r *Retrieve) Limit(limit int) *Retrieve {
	NewLimitOpt(r.Pack(), limit)
	return r

}

// || FIELDS ||

// Fields retrieves only the fields specified.
func (r *Retrieve) Fields(flds ...string) *Retrieve {
	NewFieldsOpt(r.Pack(), flds...)
	return r
}

// Relation retrieves the specified fields from the relation.
func (r *Retrieve) Relation(rel string, flds ...string) *Retrieve {
	NewRelationOpt(r.Pack(), rel, flds...)
	return r
}

// || EXEC ||

// BindExec binds Execute that Retrieve will use to run the query.
// This method MUST be called before calling Exec.
func (r *Retrieve) BindExec(e Execute) *Retrieve {
	r.Base.BindExec(e)
	return r
}

// |||| OPTS ||||

// || RELATION ||

type RelationOpt struct {
	Rel    string
	Fields FieldsOpt
}

func NewRelationOpt(p *Pack, rel string, flds ...string) {
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

// || ORDER ||

type Order int

const (
	OrderASC Order = iota + 1
	OrderDSC
)

type OrderOpt struct {
	Field string
	Order Order
}

func NewOrderOpt(p *Pack, order Order, fld string) {
	p.opts[orderOptKey] = OrderOpt{Field: fld, Order: order}
}

func RetrieveOrderOpt(p *Pack) (OrderOpt, bool) {
	qo, ok := p.opts[orderOptKey]
	if !ok {
		return OrderOpt{}, false
	}
	return qo.(OrderOpt), true
}

// || LIMIT ||

func NewLimitOpt(p *Pack, limit int) {
	p.opts[limitOptKey] = limit
}

func LimitOpt(p *Pack) (int, bool) {
	qo, ok := p.opts[limitOptKey]
	if !ok {
		return 0, false
	}
	return qo.(int), true
}
