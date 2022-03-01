package query

import "github.com/arya-analytics/aryacore/pkg/util/model"

// |||| QUERY ||||

type where struct {
	base
}

func (w *where) whereFields(flds WhereFields) {
	newWhereFieldsOpt(w.Pack(), flds)
}

func (w *where) wherePK(pk interface{}) {
	newPKOpt(w.Pack(), pk)
}

func (w *where) wherePKs(pks interface{}) {
	newPKsOpt(w.Pack(), pks)
}

// |||| WHERE FIELDS ||||

type WhereFields map[string]interface{}

type FieldOp int

//go:generate stringer -type=FieldExpOp
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
	qo := pkOpt{model.NewPKChain([]interface{}{pk})}
	p.opts[pkOptKey] = qo
}

func newPKsOpt(p *Pack, pks interface{}) {
	qo := pkOpt{model.NewPKChain(pks)}
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