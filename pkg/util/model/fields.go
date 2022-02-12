package model

import (
	"reflect"
)

type Fields struct {
	t      reflect.Type
	values []reflect.Value
}

func (f *Fields) AllNonZero() (nonZero bool) {
	for _, fld := range f.values {
		if fld.IsZero() {
			nonZero = false
		}
	}
	return nonZero
}

func (f *Fields) Raw() (rawFlds []interface{}) {
	for _, f := range f.values {
		rawFlds = append(rawFlds, f.Interface())
	}
	return rawFlds
}

func (f *Fields) ToReflect() *Reflect {
	rfl := NewReflect(reflect.New(f.t.Elem()).Interface()).NewChain()
	for _, v := range f.values {
		vRfl := newRflNilOrNonPointer(v.Interface())
		if vRfl.Type() != rfl.Type() {
			panic("incorrect type")
		}
		rfl.ChainAppend(vRfl)
	}
	return rfl
}
