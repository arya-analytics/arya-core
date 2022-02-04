package storage

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
	"strings"
)

// |||| ADAPTER ||||

// ModelExchange is a utility for
type ModelExchange struct {
	Source *model.Reflect
	Dest   *model.Reflect
}

func NewModelExchange(sourcePtr, destPtr interface{}) *ModelExchange {
	sRfl, dRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	if sRfl.RawType().Kind() != dRfl.RawType().Kind() {
		panic("model exchange received model and chain. " +
			"source and dest have same kind.")
	}
	return &ModelExchange{sRfl, dRfl}
}

func (m *ModelExchange) ToSource() {
	m.exchange(m.Dest, m.Source)
}

func (m *ModelExchange) ToDest() {
	m.exchange(m.Source, m.Dest)
}

func (m *ModelExchange) exchange(fromRfl, toRfl *model.Reflect) {
	fromRfl.ForEach(func(nFromRfl *model.Reflect, i int) {
		nToRfl := toRfl
		if toRfl.IsChain() {
			nToRfl = toRfl.ChainValueByIndexOrNew(i)
		}
		bindToSource(nToRfl, nFromRfl)
	})
}

// |||| CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(sourcePtr interface{}) interface{} {
	sourceRfl := model.NewReflect(sourcePtr)
	newRfl, ok := mc.retrieveCM(sourceRfl.Type())
	if !ok {
		panic(fmt.Sprintf("model %s could not be found in catalog", sourceRfl.Type().Name()))
	}
	if sourceRfl.IsChain() {
		return newRfl.NewChain().Pointer()
	}
	return newRfl.NewStruct().Pointer()
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

// |||| BINDING UTILITIES ||||

func bindToSource(sourceRfl, destRfl *model.Reflect) {
	for i := 0; i < destRfl.StructValue().NumField(); i++ {
		fldName, v := destRfl.Type().Field(i).Name, destRfl.StructValue().Field(i)
		fld := sourceRfl.StructValue().FieldByName(fldName)
		if validField(fld) && validValue(v) {
			if v.Type() != fld.Type() {
				if fld.Type().Kind() != reflect.Interface {
					v = exchangeNested(fld, v)
				} else if !v.Type().Implements(fld.Type()) {
					panic("doesn't implement interface")
				}
			}
			fld.Set(v)
		}
	}
}

func exchangeNested(fld, modelValue reflect.Value) reflect.Value {
	fldRfl, vRfl := newValidatedRfl(fld.Interface()), newValidatedRfl(modelValue.Interface())
	NewModelExchange(vRfl.Pointer(), fldRfl.Pointer()).ToDest()
	// If our model is a chain (i.e a slice),
	// we want to get the slice itself, not the pointer to the slice.
	if fldRfl.IsChain() {
		return fldRfl.RawValue()
	}
	return fldRfl.PointerValue()
}

func newValidatedRfl(v interface{}) *model.Reflect {
	rfl := model.UnsafeNewReflect(v)
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

func validField(fld reflect.Value) bool {
	return fld.IsValid()
}

func validValue(val reflect.Value) bool {
	return val.IsValid() && !val.IsZero()
}
