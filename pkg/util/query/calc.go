package query

import "reflect"

// Calc represents a calculation to perform on items in a data store.
type Calc int

//go:generate stringer -type=Calc
const (
	// CalcSum calculates the sum of items.
	CalcSum Calc = iota
	// CalcMax calculates the maximum value in a set of items.
	CalcMax
	// CalcMin calculates the minimum value in a set of items.
	CalcMin
	// CalcCount counts a set of items.
	CalcCount
	// CalcAVG averages a set of items.
	CalcAVG
)

// CalcOpt stores a Calc op to be performed a data store.
type CalcOpt struct {
	// Op is the calculation to be performed.
	Op Calc
	// Field is the model field to perform the calculation on.
	Field string
	// Into is the value to bind the calculation result Into.
	// NOT: Into must be a pointer.
	Into interface{}
}

// NewCalcOpt creates a new CalcOpt and binds it to the provided Pack p.
// Panics if into is not a pointer.
func NewCalcOpt(p *Pack, c Calc, fld string, into interface{}) {
	if reflect.TypeOf(into).Kind() != reflect.Ptr {
		panic("calc into arg must be ptr")
	}
	p.opts[calculateOptKey] = CalcOpt{Op: c, Field: fld, Into: into}
}

// RetrieveCalcOpt retrieves the CalcOpt from the provided Pack. Returns ok:false
// if no calculation was specified on the query.
func RetrieveCalcOpt(p *Pack) (CalcOpt, bool) {
	qo, ok := p.opts[calculateOptKey]
	if !ok {
		return CalcOpt{}, false
	}
	return qo.(CalcOpt), true
}
