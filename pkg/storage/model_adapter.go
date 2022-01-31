package storage

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

// |||| MODEL CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(sourcePtr interface{}) interface{} {
	sourceRfl := model.NewReflect(sourcePtr)
	destRfl, ok := mc.retrieveCM(sourceRfl.Type())
	if !ok {
		panic(fmt.Sprintf("model %s could not be found in catalog", sourceRfl.Type().Name()))
	}
	if sourceRfl.IsChain() {
		return destRfl.NewChain().Pointer()
	}
	return destRfl.NewStruct().Pointer()
}

func (mc ModelCatalog) NewFromType(t reflect.Type, chain bool) interface{} {
	destRfl, ok := mc.retrieveCM(t)
	if !ok {
		panic(fmt.Sprintf("model %s could not be found in catalog", t.Name()))
	}
	if chain {
		return destRfl.NewChain().Pointer()
	}
	return destRfl.NewStruct().Pointer()
}

func (mc ModelCatalog) Contains(sourcePtr interface{}) bool {
	_, ok := mc.retrieveCM(model.NewReflect(sourcePtr).Type())
	return ok
}

func (mc ModelCatalog) retrieveCM(t reflect.Type) (*model.Reflect, bool) {
	for _, destOpt := range mc {
		destOptRfl := model.NewReflect(destOpt)
		destOptRfl.Validate()
		if t.Name() == destOptRfl.Type().Name() {
			return destOptRfl, true
		}
	}
	return nil, false
}

// |||| MODEL ADAPTER ||||

type ModelAdapter struct {
	sourceRfl *model.Reflect
	destRfl   *model.Reflect
}

func NewModelAdapter(sourcePtr interface{}, destPtr interface{}) *ModelAdapter {
	sourceRfl, destRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	sourceRfl.Validate()
	destRfl.Validate()
	if sourceRfl.RawType().Kind() != destRfl.RawType().Kind() {
		panic("model adapter received model and chain. source and dest have same kind.")
	}
	return &ModelAdapter{sourceRfl, destRfl}
}

func (ma *ModelAdapter) Source() *model.Reflect {
	return ma.sourceRfl
}

func (ma *ModelAdapter) Dest() *model.Reflect {
	return ma.destRfl
}

func (ma *ModelAdapter) ExchangeToSource() {
	ma.exchange(ma.Source(), ma.Dest())
}

func (ma *ModelAdapter) ExchangeToDest() {
	ma.exchange(ma.Dest(), ma.Source())
}

func (ma *ModelAdapter) exchange(to *model.Reflect, from *model.Reflect) {
	from.ForEach(func(nRfl *model.Reflect, i int) {
		fromAm := &adaptedModel{rfl: nRfl}
		toRfl := to
		if i != -1 {
			if i >= to.ChainValue().Len() {
				toRfl = to.NewStruct()
				to.ChainAppend(toRfl)
			} else {
				toRfl = to.ChainValueByIndex(i)
			}
		}
		toAm := &adaptedModel{rfl: toRfl}
		toAm.bindVals(fromAm.mapVals())
	})
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	rfl *model.Reflect
}

type modelValues map[string]interface{}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) {
	for key, rv := range mv {
		fld := mw.rfl.StructValue().FieldByName(key)
		v := reflect.ValueOf(rv)
		if !v.IsValid() || v.IsZero() || !fld.IsValid() {
			continue
		}
		if v.Type() != fld.Type() {
			if fld.Type().Kind() == reflect.Interface {
				impl := v.Type().Implements(fld.Type())
				if !impl {
					panic("doesn't implement interface")
				}
			} else {
				fldRfl := mw.newValidatedRfl(fld.Interface())
				fldPtr := fldRfl.Pointer()
				if fldRfl.IsStruct() && fld.IsNil() {
					fldPtr = fldRfl.NewStruct().Pointer()
				}
				vPtr := mw.newValidatedRfl(rv).Pointer()
				ma := NewModelAdapter(vPtr, fldPtr)
				ma.ExchangeToDest()
				// If our model is a chain (i.e a slice),
				// we want to get the slice itself, not the pointer to the slice.
				if ma.Dest().IsChain() {
					v = ma.Dest().RawValue()
				} else {
					v = ma.Dest().PointerValue()
				}
			}
		}
		fld.Set(v)
	}
}

func (mw *adaptedModel) newValidatedRfl(v interface{}) *model.Reflect {
	rfl := model.NewReflect(v)
	// If v isn't a pointer, we need to create a pointer to it,
	// so we can manipulate its values. This is always necessary with slice fields.
	if !rfl.IsPointer() {
		rfl = rfl.ToNewPointer()
	}
	rfl.Validate()
	return rfl
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	mv := modelValues{}
	for i := 0; i < mw.rfl.StructValue().NumField(); i++ {
		fv, t := mw.rfl.StructValue().Field(i).Interface(), mw.rfl.Type().Field(i)
		mv[t.Name] = fv
	}
	return mv
}
