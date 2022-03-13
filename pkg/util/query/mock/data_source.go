package mock

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type DataSourceMem struct {
	Data model.DataSource
}

func (s *DataSourceMem) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		Create:   s.create,
		Retrieve: s.retrieve,
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
	switch calc.Op {
	case query.CalcSum:
		s.calcSum(sRfl, calc.Field, calc.Into)
	default:
		panic(fmt.Sprintf("unsupported operation %s", calc.Op))
	}
}

func (s *DataSourceMem) retrieveRelation(sRfl *model.Reflect, rel query.RelationOpt) {
	fldT := sRfl.FieldTypeByName(rel.Rel)
	if fldT.Kind() == reflect.Ptr {
		fldT = fldT.Elem()
	}
	dRfl := s.Data.Retrieve(fldT)
	_, ln := model.SplitLastFieldName(rel.Rel)
	st, ok := sRfl.StructTagChain().RetrieveByFieldName(ln)
	if !ok {
		panic(fmt.Sprintf("field %s couldn't be found on model", rel.Fields))
	}
	joinStr, ok := st.Retrieve("model", "join")
	if !ok {
		panic("model must have a join tag specified in order to perform lookup")
	}
	str := strings.Split(joinStr, "=")
	if len(str) != 2 {
		panic(fmt.Sprintf("struct tag join improperlly formatted: %s", joinStr))
	}
	dRfl.ForEach(func(nDRfl *model.Reflect, i int) {
		dFld := nDRfl.StructFieldByName(str[1])
		sRfl.ForEach(func(nSRfl *model.Reflect, i int) {
			sFld := nSRfl.StructFieldByName(str[0])
			if sFld.Interface() == dFld.Interface() {
				nSRfl.StructFieldByName(rel.Rel).Set(nDRfl.PointerValue())
			}
		})
	})
}

func (s *DataSourceMem) calcSum(sRfl *model.Reflect, field string, into interface{}) {
	intoRfl := reflect.ValueOf(into)
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		fld := rfl.StructFieldByName(field)
		if !fld.CanFloat() {
			panic("cant run a calculation on a non float!")
		}
		if intoRfl.CanFloat() {
			fldFloat := fld.Float()
			intoRflFloat := intoRfl.Float()
			intoRflFloat += fldFloat
			intoRfl.Set(reflect.ValueOf(intoRflFloat))
		}
		if intoRfl.CanInt() {
			fldInt := fld.Int()
			intoRflInt := intoRfl.Int()
			intoRflInt += fldInt
			intoRfl.Set(reflect.ValueOf(fldInt))
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
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		match := true
		for k, v := range wFld {
			fn, ln := model.SplitLastFieldName(k)
			if fn != "" {
				log.Info(rfl.FieldTypeByName(k))
				dRfl := s.Data.Retrieve(rfl.FieldTypeByName(k))
				fDrl := s.filterByWhereFields(dRfl, query.WhereFields{ln: v})
				if fDrl.ChainValue().Len() == 0 {
					match = false
				}
			} else if !fieldExpMatch(k, v, rfl) {
				match = false
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
