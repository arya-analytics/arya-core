package query

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

// |||| QUERY ||||

type Where struct {
	Base
}

func (w *Where) WhereFields(flds WhereFields) {
	NewWhereFieldsOpt(w.Pack(), flds)
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
// operations, see the FieldOp type.
type WhereFields map[string]interface{}

type FieldOp int

//go:generate stringer -type=FieldOp
const (
	FieldOpGreaterThan FieldOp = iota
	FieldOpLessThan
	FieldOpInRange
	FieldOpIn
)

type FieldExp struct {
	Op     FieldOp
	Values []interface{}
}

// GreaterThan is an option for WhereFields that will filter the results of a query by a field greater than the
// given value.
func GreaterThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldOpGreaterThan, Values: []interface{}{value}}
}

// LessThan is an option for WhereFields that will filter the results of a query by a field less than the given value.
func LessThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldOpLessThan, Values: []interface{}{value}}
}

// InRange is an option for WhereFields that will filter the results of a query by a field in the given range.
func InRange(start interface{}, stop interface{}) FieldExp {
	return FieldExp{Op: FieldOpInRange, Values: []interface{}{start, stop}}
}

// In is an option for WhereFields that will filter the results of a query by a field in the given values.
func In(values ...interface{}) FieldExp {
	return FieldExp{Op: FieldOpIn, Values: values}
}

// |||| OPTS ||||

// PKOpt is an option that will filter the results of a query by a primary key.
func PKOpt(p *Pack) (model.PKChain, bool) {
	qo, ok := p.opts[pkOptKey]
	if !ok {
		return model.PKChain{}, false
	}
	return qo.(pkOpt).PKChain, true
}

type pkOpt struct {
	PKChain model.PKChain
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
	p.SetOpt(pkOptKey, pkOpt{pkc})
}

// || WHERE FIELDS ||

// WhereFieldsOpt retrieves a WhereFields option from a query.
func WhereFieldsOpt(p *Pack) (WhereFields, bool) {
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
