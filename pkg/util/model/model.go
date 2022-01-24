package model

import (
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
)

type Reflect struct {
	modelPtr interface{}
}

func NewReflect(modelPtr interface{}) *Reflect {
	return &Reflect{
		modelPtr: modelPtr,
	}
}

func (r *Reflect) Validate() error {
	return validator.Exec(r)
}

func (r *Reflect) IsPointer() bool {
	return r.PointerType().Kind() == reflect.Ptr
}

func (r *Reflect) Pointer() interface{} {
	return r.modelPtr
}

func (r *Reflect) Type() reflect.Type {
	if r.IsChain() {
		/* raw type is the slice
		first elem is pointer to struct
		second elem is struct */
		return r.RawType().Elem().Elem()
	}
	/* raw type is pointer to struct
	first elem is struct */
	return r.RawType()
}

func (r *Reflect) Value() reflect.Value {
	if r.IsChain() {
		panic("model is a chain, cannot get a struct value")
	}
	return r.PointerValue().Elem()
}

func (r *Reflect) IsChain() bool {
	return r.RawType().Kind() == reflect.Slice
}

func (r *Reflect) IsStruct() bool {
	return r.RawType().Kind() == reflect.Struct
}

// || CHAIN METHODS ||

func (r *Reflect) ChainValue() reflect.Value {
	if r.IsChain() {
		return r.RawValue()
	}
	panic("model is a struct, cannot get a chain value")
}

func (r *Reflect) ChainAppend(v *Reflect) {
	r.ChainValue().Set(reflect.Append(r.ChainValue(), v.PointerValue()))
}

func (r *Reflect) ChainValueByIndex(i int) *Reflect {
	return NewReflect(r.ChainValue().Index(i).Interface())
}

type ForEachFunc func(rfl *Reflect, i int)

func (r *Reflect) ForEach(fef ForEachFunc) {
	if r.IsStruct() {
		fef(r, -1)
	} else {
		for i := 0; i < r.ChainValue().Len(); i++ {
			rfl := r.ChainValueByIndex(i)
			fef(rfl, i)
		}
	}
}

// || CONSTRUCTOR ||

func (r *Reflect) NewRaw() *Reflect {
	if r.IsChain() {
		return r.NewChain()
	}
	return r.NewModel()
}

func (r *Reflect) NewModel() *Reflect {
	return NewReflect(reflect.New(r.Type()).Interface())
}

func (r *Reflect) NewChain() *Reflect {
	ns := reflect.MakeSlice(reflect.SliceOf(r.NewModel().PointerType()), 0, 0)
	p := reflect.New(ns.Type())
	p.Elem().Set(ns)
	return NewReflect(p.Interface())
}

func (r *Reflect) NewPointer() *Reflect {
	p := reflect.New(r.PointerType())
	p.Elem().Set(r.PointerValue())
	return NewReflect(p.Interface())
}

// || INTERNAL TYPE + VALUE ACCESSORS ||

func (r *Reflect) PointerType() reflect.Type {
	return reflect.TypeOf(r.modelPtr)
}

func (r *Reflect) PointerValue() reflect.Value {
	return reflect.ValueOf(r.modelPtr)
}

func (r *Reflect) RawType() reflect.Type {
	return r.PointerType().Elem()
}

func (r *Reflect) RawValue() reflect.Value {
	return r.PointerValue().Elem()
}

func (r *Reflect) ValueForSet() reflect.Value {
	if r.IsChain() {
		return r.RawValue()
	}
	return r.PointerValue()
}

func validateSliceOrStruct(v interface{}) error {
	r := v.(*Reflect)
	if !r.IsStruct() && !r.IsChain() {
		return NewError(ErrTypeNonStructOrSlice)
	}
	return nil
}

func validateContainerIsPointer(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerType().Kind() != reflect.Ptr {
		return NewError(ErrTypeNonPointer)
	}
	return nil
}

var validator = validate.New([]validate.ValidateFunc{
	validateContainerIsPointer,
	validateSliceOrStruct,
})
