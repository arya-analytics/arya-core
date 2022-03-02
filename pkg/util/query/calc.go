package query

import "reflect"

type Calc int

const (
	CalcSum Calc = iota
	CalcMax
	CalcMin
	CalcCount
	CalcAVG
)

type CalcOpt struct {
	Op    Calc
	Field string
	Into  interface{}
}

func NewCalcOpt(p *Pack, c Calc, fld string, into interface{}) {
	if reflect.TypeOf(into).Kind() != reflect.Ptr {
		panic("calc into arg must be ptr")
	}
	p.opts[calculateOptKey] = CalcOpt{Op: c, Field: fld, Into: into}
}

func RetrieveCalcOpt(p *Pack) (CalcOpt, bool) {
	qo, ok := p.opts[calculateOptKey]
	if !ok {
		return CalcOpt{}, false
	}
	return qo.(CalcOpt), true
}
