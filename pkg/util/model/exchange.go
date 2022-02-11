package model

import (
	"reflect"
)

// |||| ADAPTER ||||

type FieldHandler func(sourceFldST, destFldST StructTag, sourceFld, destFld reflect.Value) (cDestFld reflect.Value, ok bool)

type Exchange struct {
	Source        *Reflect
	Dest          *Reflect
	FieldHandlers []FieldHandler
}

func NewExchange(sourcePtr, destPtr interface{}, handlers ...FieldHandler) *Exchange {
	sRfl, dRfl := NewReflect(sourcePtr), NewReflect(destPtr)
	if sRfl.RawType().Kind() != dRfl.RawType().Kind() {
		panic("model exchange received model and chain. " +
			"source and dest have same kind.")
	}
	return &Exchange{sRfl, dRfl, handlers}
}

func (m *Exchange) ToSource() {
	m.exchange(m.Dest, m.Source)
}

func (m *Exchange) ToDest() {
	m.exchange(m.Source, m.Dest)
}

func (m *Exchange) exchange(fromRfl, toRfl *Reflect) {
	fromRfl.ForEach(func(nFromRfl *Reflect, i int) {
		nToRfl := toRfl
		if toRfl.IsChain() {
			nToRfl = toRfl.ChainValueByIndexOrNew(i)
		}
		m.bindToSource(nToRfl, nFromRfl)
	})
}

// |||| BINDING UTILITIES ||||

func (m *Exchange) bindToSource(sourceRfl, destRfl *Reflect) {
	for i := 0; i < destRfl.StructValue().NumField(); i++ {
		fldName, destFld := destRfl.Type().Field(i).Name, destRfl.StructValue().Field(i)
		sourceFld := sourceRfl.StructValue().FieldByName(fldName)
		if validSourceField(sourceFld) && validDestField(destFld) {
			if destFld.Type() != sourceFld.Type() {
				if sourceFld.Type().Kind() == destFld.Type().Kind() {
					destFld = exchangeNested(sourceFld, destFld)
				} else if sourceFld.Type().Kind() != reflect.Interface {
					var ok bool
					destFld, ok = m.execCustomHandlers(fldName, sourceRfl, destRfl, sourceFld, destFld)
					if !ok {
						panic("field types incompatible")
					}
				} else if !destFld.Type().Implements(sourceFld.Type()) {
					panic("doesn't implement interface")
				}
			}
			sourceFld.Set(destFld)
		}
	}
}

const badFieldNameMsg = "couldn't find field name in struct. this is bad, and shouldn't be happening!"

func (m *Exchange) execCustomHandlers(fldName string, sourceRfl, destRfl *Reflect, sourceFld, destFld reflect.Value) (v reflect.Value, ok bool) {
	sourceTag, ok := sourceRfl.StructTagChain().RetrieveByFieldName(fldName)
	if !ok {
		panic(badFieldNameMsg)
	}
	destTag, ok := destRfl.StructTagChain().RetrieveByFieldName(fldName)
	if !ok {
		panic(badFieldNameMsg)
	}
	for _, fh := range m.FieldHandlers {
		v, ok = fh(sourceTag, destTag, sourceFld, destFld)
		if ok {
			return v, ok
		}
	}
	return v, ok
}

func exchangeNested(fld, modelValue reflect.Value) reflect.Value {
	fldRfl, vRfl := newRflNilOrNonPointer(fld.Interface()), newRflNilOrNonPointer(modelValue.Interface())
	NewExchange(vRfl.Pointer(), fldRfl.Pointer()).ToDest()
	// If our model is a chain (i.e a slice),
	// we want to get the slice itself, not the pointer to the slice.
	if fldRfl.IsChain() {
		return fldRfl.RawValue()
	}
	return fldRfl.PointerValue()
}

func newRflNilOrNonPointer(v interface{}) *Reflect {
	rfl := UnsafeNewReflect(v)
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

func validSourceField(fld reflect.Value) bool {
	return fld.IsValid()
}

func validDestField(val reflect.Value) bool {
	return val.IsValid() && !val.IsZero()
}

// |||| BASIC FIELD HANDLERS ||||

// FieldHandlerPK exchanges values between different types of primary key fields.
// To use, provide it as a field handler  a new model exchange:
//
// 		model.NewExchange(sourcePtr, destPtr, FieldHandlerPK)
//
//
// Then, the Exchange will automatically adapt PK types that wouldn't otherwise be compatible.
func FieldHandlerPK(sourceST, destST StructTag, sourceFld, destFld reflect.Value) (reflect.Value, bool) {
	if !isPKStructTag(sourceST) || !isPKStructTag(destST) {
		return reflect.Value{}, false
	}
	sourcePK, destPK := NewPK(sourceFld.Interface()), NewPK(destFld.Interface())
	oPK, err := sourcePK.NewFromString(destPK.String())
	if err != nil {
		return oPK.Value(), false
	}
	return oPK.Value(), true
}
