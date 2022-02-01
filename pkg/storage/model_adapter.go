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

func (ma *ModelAdapter) exchange(to, from *model.Reflect) {
	from.ForEach(func(fromRfl *model.Reflect, i int) {
		toRfl := to
		if to.IsChain() {
			toRfl = to.ChainValueByIndexOrNew(i)
		}
		bindToSource(toRfl, fromRfl)
	})
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

// |||| BINDING UTILITIES ||||

func bindToSource(sourceRfl, destRfl *model.Reflect) {
	for i := 0; i < destRfl.StructValue().NumField(); i++ {
		fldName, v := destRfl.Type().Field(i).Name, destRfl.StructValue().Field(i)
		fld := sourceRfl.StructValue().FieldByName(fldName)
		if validField(fld) && validValue(v) {
			if v.Type() != fld.Type() {
				if fld.Type().Kind() != reflect.Interface {
					v = adaptNestedModel(fld, v)
				} else if !v.Type().Implements(fld.Type()) {
					panic("doesn't implement interface")
				}
			}
			fld.Set(v)
		}
	}
}

func adaptNestedModel(fld reflect.Value, modelValue reflect.Value) (rv reflect.Value) {
	fldPtr := newValidatedRfl(fld.Interface()).Pointer()
	vPtr := newValidatedRfl(modelValue.Interface()).Pointer()
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

func newValidatedRfl(v interface{}) *model.Reflect {
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

func validField(fld reflect.Value) bool {
	return fld.IsValid()
}

func validValue(val reflect.Value) bool {
	return val.IsValid() && !val.IsZero()
}
