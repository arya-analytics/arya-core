package model

import (
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

type WhereFields map[string]interface{}
