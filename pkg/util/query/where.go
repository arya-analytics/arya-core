package query

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

// |||| QUERY ||||

type where struct {
	base
}

func (w *where) whereFields(flds WhereFields) {
	newWhereFieldsOpt(w.Pack(), flds)
}

func (w *where) wherePK(pk interface{}) {
	if reflect.TypeOf(pk).Kind() == reflect.Slice {
		panic("wherepk can't be called with multiple primary keys!")
	}
	newPKOpt(w.Pack(), pk)
}

func (w *where) wherePKs(pks interface{}) {
	if reflect.TypeOf(pks).Kind() != reflect.Slice {
		panic("wherepks can't be called with a single primary key!")
	}
	newPKOpt(w.Pack(), pks)
}

// |||| WHERE FIELDS ||||

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
	Op   FieldOp
	Vals []interface{}
}

func GreaterThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldOpGreaterThan, Vals: []interface{}{value}}
}

func LessThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldOpLessThan, Vals: []interface{}{value}}
}

func InRange(start interface{}, stop interface{}) FieldExp {
	return FieldExp{Op: FieldOpInRange, Vals: []interface{}{start, stop}}
}

func In(vals ...interface{}) FieldExp {
	return FieldExp{Op: FieldOpIn, Vals: vals}
}

// |||| OPTS ||||

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

func newPKOpt(p *Pack, pk interface{}) {
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

	qo := pkOpt{pkc}
	p.opts[pkOptKey] = qo
}

// || WHERE FIELDS ||

func WhereFieldsOpt(p *Pack) (WhereFields, bool) {
	qo, ok := p.opts[whereFieldsOptKey]
	if !ok {
		return WhereFields{}, false
	}
	return qo.(WhereFields), true
}

func newWhereFieldsOpt(p *Pack, flds WhereFields) {
	p.opts[whereFieldsOptKey] = flds
}
