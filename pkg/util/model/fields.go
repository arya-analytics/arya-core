package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Fields struct {
	source  *Reflect
	fldName string
}

func (f *Fields) values() (values []reflect.Value) {
	f.source.ForEach(func(rfl *Reflect, i int) {
		values = append(values, rfl.StructFieldByName(f.fldName))
	})
	return values
}

func (f *Fields) AllNonZero() bool {
	allNonZero := true
	for _, fld := range f.values() {
		if !fld.IsZero() {
			allNonZero = false
		}
	}
	return allNonZero
}

func (f *Fields) Raw() (rawFlds []interface{}) {
	for _, v := range f.values() {
		rawFlds = append(rawFlds, v.Interface())
	}
	return rawFlds
}

func (f *Fields) ToPKChain() PKChain {
	return NewPKChain(f.Raw())
}

func (f *Fields) ToReflect() *Reflect {
	t, ok := f.source.Type().FieldByName(f.fldName)
	if !ok {
		panic(fmt.Sprintf("field %s does not exist", f.fldName))
	}
	newRfl := NewReflect(reflect.New(t.Type.Elem()).Interface()).NewChain()
	f.source.ForEach(func(rfl *Reflect, i int) {
		fld := rfl.StructFieldByName(f.fldName)
		fldRfl := newReflectNilNonPtr(fld.Interface())
		fld.Set(fldRfl.PointerValue())
		newRfl.ChainAppend(fldRfl)
	})
	return newRfl
}

func SplitFieldNames(name string) []string {
	return strings.Split(name, ".")
}

func SplitLastFieldName(name string) (string, string) {
	sn := SplitFieldNames(name)
	fn := strings.Join(sn[0:len(sn)-1], ".")
	return fn, sn[len(sn)-1]
}

func SplitFirstFieldName(name string) string {
	return SplitFieldNames(name)[0]
}

type FieldExpOp int


//go:generate stringer -type=FieldExpOp
const (
	FieldExpOpGreaterThan FieldExpOp = iota
	FieldExpOpLessThan
	FieldExpOpInRange
)

type FieldExp struct {
	Op   FieldExpOp
	Vals []interface{}
}

func FieldGreaterThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldExpOpGreaterThan, Vals: []interface{}{value}}
}

func FieldLessThan(value interface{}) FieldExp {
	return FieldExp{Op: FieldExpOpLessThan, Vals: []interface{}{value}}
}

func FieldInRange(start interface{}, stop interface{}) FieldExp {
	return FieldExp{Op: FieldExpOpInRange, Vals: []interface{}{start, stop}}
}

// What are the things we need to return
// 1. What type of expression is this
// 2. What are the parameters of the expression

type WhereFields map[string]interface{}
