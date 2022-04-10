// Package filter filters an input list of items based on a provided query.
// To use the filter, call Filter.
//
// Example:
//      inputData := []*models.Item{{Name: "foo"}, {Name: "bar"}}
// 		outputData := filter.Filter(query.NewRetrieve().WhereFields(query.WhereFields{"name": "foo"}).Pack(), inputData)
//      fmt.Println(outputData)
// 	 	// [{Name:foo}]
//
// Also provides filters for individual options in a query:
//
// 		PK - By primary key (PK)
//      WhereFields - By where fields. Assumes that relations are already defined on the input data. Does not support
//      	(field expressions) query.FieldExpression such as query.InRange.
// 		Calc - running calculations (Calc) - Only supports addition.
//
package filter

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"reflect"
)

// Filter filters a slice of models based on the provided query.
// It supports the following query options as filter parameters:
//
//      - by primary key (PK)
//      - by where fields (WhereFields) - Assumes that relations are already defined on the input data.
// 		- running calculations (Calc) - Only supports addition.
//
// The provided data must be either a pointer to a slice of structs or a *model.Reflect pointing to a slice of structs.
// Returns an output reflection as pointer to the slice of filtered values. Also returns any errors encountered.
//
// Options:
//		ErrorOnNotFound: If true, an error is returned if no results are found.
//      ErrorOnMultiple: If true, an error is returned if multiple results are found.
//
// NOTE: Although the output data is a new slice, the values inside may still maintain a reference to the original data.
func Filter[T any](p *query.Pack, data []T, opts ...Opt) ([]T, error) {
	res, err := ReflectFilter(p, model.NewReflect(&data), opts...)
	return res.ChainValue().Interface().([]T), err
}

// ReflectFilter is the same as Filter, but takes a *model.Reflect as input.
func ReflectFilter(p *query.Pack, d *model.Reflect, opts ...Opt) (*model.Reflect, error) {
	if !d.IsChain() {
		panic("filter: model is not chain")
	}
	fo := newOpts(opts...)
	var oRfl *model.Reflect
	for i, f := range filters() {
		if i == 0 {
			oRfl = f(p, d)
		} else {
			oRfl = f(p, oRfl)
		}
	}
	return oRfl, optErr(fo, oRfl)
}

type Opt func(fo *fOpts)

// ErrorOnNotFound returns an error if no results are found.
func ErrorOnNotFound() Opt {
	return func(fo *fOpts) {
		fo.errOnNotFound = true
	}
}

// ErrorOnMultiple returns an error if multiple results are found.
func ErrorOnMultiple() Opt {
	return func(fo *fOpts) {
		fo.errOnMultiple = true
	}
}

// PK filters the input data by the primary key option of the query.
// Returns the original data if no primary key was specified on the query.
func PK(p *query.Pack, sRfl *model.Reflect) *model.Reflect {
	pkc, ok := query.RetrievePKOpt(p)
	if !ok {
		return sRfl
	}
	return sRfl.Filter(func(rfl *model.Reflect, _ int) bool {
		pk := rfl.PK()
		for _, o := range pkc {
			if pk.Equals(o) {
				return true
			}
		}
		return false
	})
}

// WhereFields filters the input data by the where fields option of the query.
// Returns the original data if no where fields were specified on the query.
func WhereFields(p *query.Pack, sRfl *model.Reflect) *model.Reflect {
	wf, ok := query.RetrieveWhereFieldsOpt(p)
	if !ok {
		return sRfl
	}
	return sRfl.Filter(func(rfl *model.Reflect, _ int) bool {
		match := true
		for k, v := range wf {
			if !fieldsMatch(k, v, rfl) {
				match = false
			}
		}
		return match
	})
}

// |||| CALCULATE ||||

// Calc executes calculations on the input data by the calc option of the query.
// Returns the original data. Calculations are only supported for addition.
func Calc(p *query.Pack, sRfl *model.Reflect) *model.Reflect {
	calc, ok := query.RetrieveCalcOpt(p)
	if !ok {
		return sRfl
	}
	reflect.ValueOf(calc.Into).Elem().Set(reflect.Zero(reflect.TypeOf(calc.Into).Elem()))
	switch calc.Op {
	case query.CalcSum:
		calcSum(sRfl, calc.Field, calc.Into)
	default:
		panic(fmt.Sprintf("unsupported operation %s", calc.Op))
	}
	return sRfl
}

func calcSum(sRfl *model.Reflect, field string, into interface{}) {
	intoRfl := reflect.ValueOf(into)
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		fld := rfl.StructFieldByName(field)
		if !fld.CanFloat() && !fld.CanInt() {
			panic("cant run a calculation on a non number!")
		}
		if fld.CanFloat() {
			fldFloat := fld.Float()
			intoRflFloat := intoRfl.Elem().Float()
			intoRflFloat += fldFloat
			intoRfl.Elem().Set(reflect.ValueOf(intoRflFloat))
		}
		if fld.CanInt() {
			fldInt := fld.Int()
			intoRflInt := intoRfl.Elem().Int()
			intoRflInt += fldInt
			intoRfl.Elem().Set(reflect.ValueOf(intoRflInt))
		}
	})
}

// |||| FILTER COMPOSITION ||||

type filter func(p *query.Pack, rfl *model.Reflect) *model.Reflect

func filters() []filter {
	return []filter{
		PK,
		WhereFields,
		Calc,
	}
}

// |||| FIELD MATCH ||||

func fieldsMatch(wFldName string, wFldVal interface{}, source *model.Reflect) bool {
	fldVal := source.StructFieldByName(wFldName)
	if !fldVal.IsValid() {
		return false
	}
	_, ok := wFldVal.(query.FieldExpression)
	// We don't currently support any field expressions.
	if ok {
		panic("field expressions not currently supported")
	}
	return wFldVal == fldVal.Interface()
}

// |||| OPTIONS ||||

type fOpts struct {
	errOnNotFound bool
	errOnMultiple bool
}

func newOpts(opts ...Opt) *fOpts {
	fo := &fOpts{}
	for _, o := range opts {
		o(fo)
	}
	return fo
}

func optErr(fo *fOpts, rfl *model.Reflect) error {
	if fo.errOnNotFound && rfl.ChainValue().Len() == 0 {
		return query.Error{
			Type:    query.ErrorTypeItemNotFound,
			Message: fmt.Sprintf("filter: no results for query %s", rfl.ChainValue()),
			Base:    nil,
		}
	}
	if fo.errOnMultiple && rfl.ChainValue().Len() > 1 {
		return query.Error{
			Type:    query.ErrorTypeMultipleResults,
			Message: fmt.Sprintf("filter: multiple results for query %s", rfl.ChainValue()),
			Base:    nil,
		}
	}
	return nil
}
