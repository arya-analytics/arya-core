package query

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

// |||| QUERY ||||

type Where struct {
	Base
}

func (w *Where) WhereFields(fields WhereFields) {
	NewWhereFieldsOpt(w.Pack(), fields)
}

func (w *Where) WherePK(pk interface{}) {
	if reflect.TypeOf(pk).Kind() == reflect.Slice {
		panic("wherepk can't be called with multiple primary keys!")
	}
	NewPKOpt(w.Pack(), pk)
}

func (w *Where) WherePKs(pks interface{}) {
	NewPKOpt(w.Pack(), pks)
}

// |||| WHERE FIELDS ||||

// WhereFields holds parameters that can be used to filter the results of a query on arbitrary fields.
//
// The key of the map is the field name, and the value is the value to filter by.
//
// The value can be wrapped in a specific matcher operation such as GreaterThan, LessThan, or InRange. For all available
// operations, see the FieldFilter type.
type WhereFields map[string]interface{}

type FieldFilter int

//go:generate stringer -type=FieldOp
const (
	FieldFilterGreaterThan FieldFilter = iota
	FieldFilterLessThan
	FieldFilterInRange
	FilterFilterIsIn
)

type FieldExpression struct {
	Op     FieldFilter
	Values []interface{}
}

// GreaterThan is an option for WhereFields that will filter the results of a query by a field greater than the
// given value.
func GreaterThan(value interface{}) FieldExpression {
	return FieldExpression{Op: FieldFilterGreaterThan, Values: []interface{}{value}}
}

// LessThan is an option for WhereFields that will filter the results of a query by a field less than the given value.
func LessThan(value interface{}) FieldExpression {
	return FieldExpression{Op: FieldFilterLessThan, Values: []interface{}{value}}
}

// InRange is an option for WhereFields that will filter the results of a query by a field in the given range.
func InRange(start interface{}, stop interface{}) FieldExpression {
	return FieldExpression{Op: FieldFilterInRange, Values: []interface{}{start, stop}}
}

// IsIn is an option for WhereFields that will filter the results of a query by a field in the given values.
func IsIn(values ...interface{}) FieldExpression {
	return FieldExpression{Op: FilterFilterIsIn, Values: values}
}

// |||| OPTS ||||

// RetrievePKOpt is an option that will filter the results of a query by a primary key.
func RetrievePKOpt(p *Pack, opts ...OptRetrieveOpt) (model.PKChain, bool) {
	qo, ok := p.RetrieveOpt(pkOptKey, opts...)
	if !ok {
		return model.PKChain{}, false
	}
	return qo.(model.PKChain), true
}

// NewPKOpt creates a new option that stores a primary key.
func NewPKOpt(p *Pack, pk interface{}) {
	var pkc model.PKChain
	switch pk.(type) {
	case model.PK:
		pkc = model.PKChain{pk.(model.PK)}
	case model.PKChain:
		pkc = pk.(model.PKChain)
	default:
		if reflect.TypeOf(pk).Kind() == reflect.Slice {
			pkc = model.NewPKChain(pk)
		} else {
			pkc = model.NewPKChain([]interface{}{pk})
		}
	}
	p.SetOpt(pkOptKey, pkc)
}

// || WHERE FIELDS ||

// RetrieveWhereFieldsOpt retrieves a WhereFields option from a query.
func RetrieveWhereFieldsOpt(p *Pack) (WhereFields, bool) {
	qo, ok := p.RetrieveOpt(whereFieldsOptKey)
	if !ok {
		return WhereFields{}, false
	}
	return qo.(WhereFields), true
}

// NewWhereFieldsOpt applies a WhereFields option to a query.
func NewWhereFieldsOpt(p *Pack, fields WhereFields) {
	p.SetOpt(whereFieldsOptKey, fields)
}
