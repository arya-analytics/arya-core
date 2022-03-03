package model

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"reflect"
)

// |||| ADAPTER ||||

type FieldHandler func(sourceFldST, destFldST StructTag, sourceFld, destFld reflect.Value) (cDestFld reflect.Value, ok bool)

type Exchange struct {
	source        *Reflect
	dest          *Reflect
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

func (m *Exchange) Source() *Reflect {
	return m.source
}

func (m *Exchange) Dest() *Reflect {
	return m.dest
}

func (m *Exchange) ToSource() {
	m.exchange(m.Dest(), m.Source())
}

func (m *Exchange) ToDest() {
	m.exchange(m.Source(), m.Dest())
}

func (m *Exchange) exchange(fromRfl, toRfl *Reflect) {
	update := exchangeIsUpdate(toRfl)
	fromRfl.ForEach(func(nFromRfl *Reflect, i int) {
		m.bindToSource(nToRfl(toRfl, nFromRfl, i, update), nFromRfl)
	})
}

/// |||| RETRIEVE UTILITIES ||||

func exchangeIsUpdate(toRfl *Reflect) bool {
	return toRfl.IsChain() && toRfl.ChainValue().Len() > 0 && !toRfl.PKChain().AllZero()
}

func warnBadUpdate(toRfl, nFromRfl *Reflect) {
	log.WithFields(log.Fields{
		"exchangingFrom": nFromRfl.Type().Name(),
		"fromPK":         nFromRfl.PK(),
		"exchangingTo":   toRfl.Type().Name(),
	}).Warn("model exchange doing update, but appears PK does not exist in dest. This may lead to strange bugs.")
}

func nToRfl(toRfl, nFromRfl *Reflect, index int, update bool) *Reflect {
	if !toRfl.IsChain() {
		return toRfl
	}
	if update {
		if !nFromRfl.PK().IsZero() {
			rfl, ok := toRfl.ValueByPK(nFromRfl.PK())
			if ok {
				return rfl
			} else {
				warnBadUpdate(toRfl, nFromRfl)
			}
		}
	}
	return toRfl.ChainValueByIndexOrNew(index)
}

// |||| BINDING UTILITIES ||||

func (m *Exchange) bindToSource(sourceRfl, destRfl *Reflect) {
	for i := 0; i < destRfl.StructValue().NumField(); i++ {
		fldName, destFld := destRfl.Type().Field(i).Name, destRfl.StructValue().Field(i)
		sourceFld := sourceRfl.StructFieldByName(fldName)
		if validSourceField(sourceFld) && validDestField(destFld) {
			if destFld.Type() != sourceFld.Type() {
				if sourceFld.Type().Kind() == destFld.Type().Kind() {
					destFld = exchangeNested(sourceFld, destFld)
				} else if sourceFld.Type().Kind() != reflect.Interface {
					var ok bool
					destFld, ok = m.execCustomHandlers(fldName, sourceRfl, destRfl, sourceFld, destFld)
					if !ok {
						panic(fmt.Sprintf("field %s has incompatible type %s", fldName, sourceFld.Type()))
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
	fldRfl, vRfl := newReflectNilNonPtr(fld.Interface()), newReflectNilNonPtr(modelValue.Interface())
	NewExchange(vRfl.Pointer(), fldRfl.Pointer()).ToDest()
	// If our model is a chain (i.e a slice),
	// we want to get the slice itself, not the pointer to the slice.
	if fldRfl.IsChain() {
		return fldRfl.RawValue()
	}
	return fldRfl.PointerValue()
}

func newReflectNilNonPtr(v interface{}) *Reflect {
	rfl := UnsafeNewReflect(v)
	// If v isn't a pointer, we need to create a pointer to it,
	// so we can manipulate its values. This is always necessary with slice fields.
	if !rfl.IsPointer() {
		rfl = rfl.ToNewPointer()
	}
	// If v is zero, that means it's a struct we can't assign values to,
	// so we need to initialize a new empty struct with a non-zero Val.
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
	if !isPKField(sourceST) || !isPKField(destST) {
		return reflect.Value{}, false
	}
	sourcePK, destPK := NewPK(sourceFld.Interface()), NewPK(destFld.Interface())
	oPK, err := sourcePK.NewFromString(destPK.String())
	if err != nil {
		return oPK.Value(), false
	}
	return oPK.Value(), true
}

func isPKField(st StructTag) bool {
	return st.Field.Type == reflect.TypeOf(uuid.UUID{}) || st.Field.Type.Kind() == reflect.String
}
