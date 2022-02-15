package model

import (
	"reflect"
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

func (f *Fields) AllNonZero() (nonZero bool) {
	for _, fld := range f.values() {
		if fld.IsZero() {
			nonZero = false
		}
	}
	return nonZero
}

func (f *Fields) Raw() (rawFlds []interface{}) {
	for _, v := range f.values() {
		rawFlds = append(rawFlds, v.Interface())
	}
	return rawFlds
}

func (f *Fields) ToReflect() *Reflect {
	t, ok := f.source.Type().FieldByName(f.fldName)
	if !ok {
		panic("field does not exist")
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
