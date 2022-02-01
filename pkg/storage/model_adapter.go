package storage

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
	"strings"
)

// |||| ADAPTER ||||

type ModelAdapter struct {
	sourceRfl *model.Reflect
	destRfl   *model.Reflect
}

func NewModelAdapter(sourcePtr interface{}, destPtr interface{}) *ModelAdapter {
	sRfl, dRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	if sRfl.RawType().Kind() != dRfl.RawType().Kind() {
		panic("model adapter received model and chain. source and dest have same kind.")
	}
	return &ModelAdapter{sRfl, dRfl}
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

// || BINDING ||

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) {
	for key, rv := range mv {
		fld := mw.rfl.StructValue().FieldByName(key)
		v := reflect.ValueOf(rv)
		if !mw.validField(fld) || !mw.validValue(v) {
			continue
		}
		if v.Type() != fld.Type() {
			if fld.Type().Kind() != reflect.Interface {
				v = mw.adaptNested(fld, v)
			} else if !v.Type().Implements(fld.Type()) {
				panic("doesn't implement interface")
			}
		}
		fld.Set(v)
	}
}

func (mw *adaptedModel) adaptNested(fld reflect.Value,
	modelValue reflect.Value) (rv reflect.Value) {
	fldPtr := mw.newValidatedRfl(fld.Interface()).Pointer()
	vPtr := mw.newValidatedRfl(modelValue.Interface()).Pointer()
	ma := NewModelAdapter(vPtr, fldPtr)
	ma.ExchangeToDest()
	// If our model is a chain (i.e a slice),
	// we want to get the slice itself, not the pointer to the slice.
	if ma.Dest().IsChain() {
		rv = ma.Dest().RawValue()
	} else {
		rv = ma.Dest().PointerValue()
	}
	return rv
}

func (mw *adaptedModel) newValidatedRfl(v interface{}) *model.Reflect {
	rfl := model.UnsafeUnvalidatedNewReflect(v)
	// If v isn't a pointer, we need to create a pointer to it,
	// so we can manipulate its values. This is always necessary with slice fields.
	if !rfl.IsPointer() {
		rfl = rfl.ToNewPointer()
	}
	// If v is zero, that means it's a struct we can't assign values to,
	// so we need to initialize a new empty struct with a non-zero value.
	if rfl.PointerValue().IsZero() {
		rfl = rfl.NewRaw()
	}
	rfl.Validate()
	return rfl
}

func (mw *adaptedModel) validField(fld reflect.Value) bool {
	return fld.IsValid()
}

func (mw *adaptedModel) validValue(val reflect.Value) bool {
	return val.IsValid() && !val.IsZero()
}

// || MAPPING ||

type modelValues map[string]interface{}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	mv := modelValues{}
	for i := 0; i < mw.rfl.StructValue().NumField(); i++ {
		fv, t := mw.rfl.StructValue().Field(i).Interface(), mw.rfl.Type().Field(i)
		mv[t.Name] = fv
	}
	return mv
}

// |||| CATALOG ||||

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

func (mc ModelCatalog) Contains(sourcePtr interface{}) bool {
	_, ok := mc.retrieveCM(model.NewReflect(sourcePtr).Type())
	return ok
}

func (mc ModelCatalog) retrieveCM(t reflect.Type) (*model.Reflect, bool) {
	for _, opt := range mc {
		optRfl := model.NewReflect(opt)
		if strings.EqualFold(t.Name(), optRfl.Type().Name()) {
			return optRfl, true
		}
	}
	return nil, false
}
