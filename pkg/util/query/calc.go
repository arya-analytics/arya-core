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
	Op      Calc
	FldName string
	Into    interface{}
}

func newCalcOpt(q *Pack, c Calc, fldName string, into interface{}) {
	if reflect.TypeOf(into).Kind() != reflect.Ptr {
		panic("calc into arg must be ptr")
	}
	q.opts[calculateOptKey] = CalcOpt{Op: c, FldName: fldName, Into: into}
}

func RetrieveCalcOpt(q *Pack) (CalcOpt, bool) {
	qo, ok := q.opts[calculateOptKey]
	if !ok {
		return CalcOpt{}, false
	}
	return qo.(CalcOpt), true
}
