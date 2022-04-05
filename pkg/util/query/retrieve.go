package query

// Retrieve retrieves a model or set of models depending on the parameters passed.
type Retrieve struct {
	Where
}

// || CONSTRUCTOR ||

// NewRetrieve opens a new Retrieve query.
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

// WherePK queries by the primary of the model to be retrieved.
func (r *Retrieve) WherePK(pk interface{}) *Retrieve {
	r.Where.WherePK(pk)
	return r
}

// WherePKs queries by a set of primary keys of models to be retrieved.
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
func (r *Retrieve) Order(dir OrderDirection, fld string) *Retrieve {
	NewOrderOpt(r.Pack(), dir, fld)
	return r
}

// || LIMIT ||

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

// || MEMO |||

func (r *Retrieve) WithMemo(memo *Memo) *Retrieve {
	NewMemoOpt(r.Pack(), memo)
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

// RelationOpt is an option for retrieving a model's relation.
type RelationOpt struct {
	// Name specifies the name of the relation to retrieve.
	Name string
	// Fields specifies the fields to retrieve from the relation.
	Fields FieldsOpt
}

// NewRelationOpt creates a new RelationOpt.
func NewRelationOpt(p *Pack, name string, fields ...string) {
	nro := RelationOpt{name, fields}
	ro, ok := p.RetrieveOpt(relationOptKey)
	if !ok {
		p.SetOpt(relationOptKey, []RelationOpt{nro})
	} else {
		p.SetOpt(relationOptKey, append(ro.([]RelationOpt), nro))
	}
}

// RelationOpts retrieves a slice of all RelationOpt applied to the query.
func RelationOpts(p *Pack) []RelationOpt {
	o, ok := p.RetrieveOpt(relationOptKey)
	if !ok {
		return []RelationOpt{}
	}
	return o.([]RelationOpt)
}

// || ORDER ||

// OrderDirection is a type that specifies the direction in which to order the results of a query.
type OrderDirection int

const (
	// OrderASC orders the results in ascending order.
	OrderASC OrderDirection = iota + 1
	// OrderDSC orders the results in descending order.
	OrderDSC
)

// OrderOpt is an option for ordering the results of a query.
type OrderOpt struct {
	// Field represents the field to order by.
	Field string
	// Direction stores the direction in which to order the results.
	Direction OrderDirection
}

// NewOrderOpt creates a new OrderOpt.
func NewOrderOpt(p *Pack, order OrderDirection, fld string) {
	p.SetOpt(orderOptKey, OrderOpt{fld, order})
}

// RetrieveOrderOpt retrieves any order options applied to the query.
// Returns false for the second argument if no ordering was specified.
func RetrieveOrderOpt(p *Pack) (OrderOpt, bool) {
	qo, ok := p.RetrieveOpt(orderOptKey)
	if !ok {
		return OrderOpt{}, false
	}
	return qo.(OrderOpt), true
}

// || LIMIT ||

// NewLimitOpt creates a new LimitOpt.
func NewLimitOpt(p *Pack, limit int) {
	p.opts[limitOptKey] = limit
}

// LimitOpt is an option for limiting the number of results returned by a query.
func LimitOpt(p *Pack) (int, bool) {
	qo, ok := p.opts[limitOptKey]
	if !ok {
		return 0, false
	}
	return qo.(int), true
}

// || MEMO ||

// NewMemoOpt creates a new MemoOpt.
func NewMemoOpt(p *Pack, memo *Memo) {
	p.SetOpt(memoOptKey, memo)
}

// MemoOpt is an option that that applies a Memo to the query.
// For more information on memoizing query results, see Memo.
func MemoOpt(p *Pack) (*Memo, bool) {
	qo, ok := p.RetrieveOpt(memoOptKey)
	if !ok {
		return nil, false
	}
	return qo.(*Memo), true
}
