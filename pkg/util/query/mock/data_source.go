package mock

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"reflect"
	"strings"
)

type DataSourceMem struct {
	query.AssembleBase
	Data model.DataSource
}

func NewDataSourceMem() *DataSourceMem {
	ds := &DataSourceMem{Data: map[reflect.Type]*model.Reflect{}}
	ds.AssembleBase = query.NewAssemble(ds.Exec)
	return ds
}

func (s *DataSourceMem) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		Create:   s.create,
		Retrieve: s.retrieve,
		Update:   s.update,
	})
}

func (s *DataSourceMem) retrieve(ctx context.Context, p *query.Pack) error {
	f := s.filter(p)
	if f.ChainValue().Len() == 0 {
		return query.NewSimpleError(query.ErrorTypeItemNotFound, nil)
	}
	var exc *model.Exchange
	if p.Model().IsStruct() {
		exc = model.NewExchange(p.Model().Pointer(), f.ChainValueByIndex(0).Pointer())
	} else {
		exc = model.NewExchange(p.Model().Pointer(), f.Pointer())
	}
	exc.ToSource()
	return nil
}

func (s *DataSourceMem) create(ctx context.Context, p *query.Pack) error {
	dRfl := s.Data.Retrieve(p.Model().Type())
	p.Model().ForEach(func(rfl *model.Reflect, i int) {
		dRfl.ChainAppendEach(rfl)
	})
	return nil
}

func (s *DataSourceMem) update(ctx context.Context, p *query.Pack) error {
	bulk := query.BulkUpdateOpt(p)
	if bulk {
		return s.bulkUpdate(ctx, p)
	}
	pk, ok := query.PKOpt(p)
	if !ok {
		panic("non-bulk update must havea pk")
	}
	dRfl := s.filterByPK(s.Data.Retrieve(p.Model().Type()), pk)
	if dRfl.ChainValue().Len() == 0 {
		return query.NewSimpleError(query.ErrorTypeItemNotFound, nil)
	}
	if dRfl.ChainValue().Len() > 1 {
		return query.NewSimpleError(query.ErrorTypeMultipleResults, nil)
	}
	fo, ok := query.RetrieveFieldsOpt(p)
	if !ok {
		panic("fields must be specified for updates")
	}
	p.Model().ForEach(func(nSRfl *model.Reflect, i int) {
		dRfl.ForEach(func(nDRfl *model.Reflect, i int) {
			for _, f := range fo {
				nDRfl.StructFieldByName(f).Set(nSRfl.StructFieldByName(f))
			}
		})
	})
	return nil
}

func (s *DataSourceMem) bulkUpdate(ctx context.Context, p *query.Pack) error {
	dRfl := s.Data.Retrieve(p.Model().Type())
	fo, ok := query.RetrieveFieldsOpt(p)
	if !ok {
		panic("fields must be specified for bulk updates")
	}
	p.Model().ForEach(func(nSRfl *model.Reflect, i int) {
		dRfl.ForEach(func(nDRfl *model.Reflect, i int) {
			if nDRfl.PK() == nSRfl.PK() {
				for _, f := range fo {
					nDRfl.StructFieldByName(f).Set(nSRfl.StructFieldByName(f))
				}
			}
		})
	})
	return nil
}

func (s *DataSourceMem) filter(p *query.Pack) *model.Reflect {
	var filteredRfl = s.Data.Retrieve(p.Model().Type())
	pkC, ok := query.PKOpt(p)
	if ok {
		filteredRfl = s.filterByPK(filteredRfl, pkC)
	}
	wFld, ok := query.WhereFieldsOpt(p)
	if ok {
		filteredRfl = s.filterByWhereFields(filteredRfl, wFld)
	}
	calcOpt, ok := query.RetrieveCalcOpt(p)
	if ok {
		s.runCalculations(filteredRfl, calcOpt)
	}
	ro := query.RelationOpts(p)
	for _, r := range ro {
		s.retrieveRelation(filteredRfl, r)
	}
	return filteredRfl
}

func (s *DataSourceMem) runCalculations(sRfl *model.Reflect, calc query.CalcOpt) {
	reflect.ValueOf(calc.Into).Elem().Set(reflect.Zero(reflect.TypeOf(calc.Into).Elem()))
	switch calc.Op {
	case query.CalcSum:
		s.calcSum(sRfl, calc.Field, calc.Into)
	default:
		panic(fmt.Sprintf("unsupported operation %s", calc.Op))
	}
}

func (s *DataSourceMem) retrieveRelation(sRfl *model.Reflect, rel query.RelationOpt) {
	names := model.SplitFieldNames(rel.Rel)
	name := names[0]
	fldT := sRfl.FieldTypeByName(name)
	st, ok := sRfl.StructTagChain().RetrieveByFieldName(name)
	if !ok {
		panic(fmt.Sprintf("field %s couldn't be found on model %s", name, sRfl.Type()))
	}
	if fldT.Kind() == reflect.Ptr {
		fldT = fldT.Elem()
	}
	joinStr, ok := st.Retrieve("model", "join")
	if !ok {
		panic("model must have a join tag specified in order to perform lookup")
	}
	str := strings.Split(joinStr, "=")
	if len(str) != 2 {
		panic(fmt.Sprintf("struct tag join improperlly formatted: %s", joinStr))
	}
	s.Data.Retrieve(fldT).ForEach(func(nDRfl *model.Reflect, i int) {
		dFld := nDRfl.StructFieldByName(str[1])
		sRfl.ForEach(func(nSRfl *model.Reflect, i int) {
			sFld := nSRfl.StructFieldByName(str[0])
			if sFld.Interface() == dFld.Interface() {
				if len(names) > 1 {
					s.retrieveRelation(nDRfl, query.RelationOpt{
						Rel:    strings.Join(names[1:], "."),
						Fields: rel.Fields,
					})
				}
				nSRfl.StructFieldByName(names[0]).Set(nDRfl.PointerValue())
			}
		})
	})
}

func (s *DataSourceMem) calcSum(sRfl *model.Reflect, field string, into interface{}) {
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

func (s *DataSourceMem) filterByPK(sRfl *model.Reflect, pkc model.PKChain) *model.Reflect {
	nRfl := sRfl.NewChain()
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		for _, pk := range pkc {
			if rfl.PK().Equals(pk) {
				nRfl.ChainAppendEach(rfl)
			}
		}
	})
	return nRfl
}

func (s *DataSourceMem) filterByWhereFields(sRfl *model.Reflect, wFld query.WhereFields) *model.Reflect {
	nRfl := sRfl.NewChain()
	for k, _ := range wFld {
		fn, ln := model.SplitLastFieldName(k)
		if fn != "" {
			s.retrieveRelation(sRfl, query.RelationOpt{Rel: fn, Fields: query.FieldsOpt{ln}})
		}
	}
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		match := false
		for k, v := range wFld {
			if fieldExpMatch(k, v, rfl) {
				match = true
			}
		}
		if match {
			nRfl.ChainAppendEach(rfl)
		}
	})
	return nRfl
}

func fieldExpMatch(wFldName string, wFldVal interface{}, source *model.Reflect) bool {
	fldVal := source.StructFieldByName(wFldName)
	if fldVal.IsZero() {
		return false
	}
	_, ok := wFldVal.(query.FieldExp)
	// We don't currently support any field expressions.
	if ok {
		panic("field expressions not currently supported")
	}
	return wFldVal == fldVal.Interface()
}
