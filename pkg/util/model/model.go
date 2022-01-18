package model

import (
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type Reflect struct {
	modelPtr interface{}
}

func NewReflect(modelPtr interface{}) (r *Reflect) {
	r = &Reflect{
		modelPtr: modelPtr,
	}
	return r
}

func (r *Reflect) Validate() error {
	return validator.Exec(r.modelPtr)
}

func (r *Reflect) Pointer() interface{} {
	return r.modelPtr
}

func (r *Reflect) IsChain() bool {
	return r.Type().Kind() == reflect.Slice
}

func (r *Reflect) IsStruct() bool {
	return r.Type().Kind() == reflect.Struct
}

func (r *Reflect) Name() string {
	if r.IsChain() {
		return r.Type().Elem().Elem().Name()
	}
	return r.Type().Name()
}

func (r *Reflect) containerType() reflect.Type {
	return reflect.TypeOf(r.modelPtr)
}

func (r *Reflect) containerValue() reflect.Value {
	return reflect.ValueOf(r.modelPtr)
}

// || CHAIN METHODS ||

func (r *Reflect) ChainValue() reflect.Value {
	if r.IsChain() {
		return r.Value()
	}
	log.Fatalln("model is not a chain, cannot extract its value")
	return reflect.Value{}
}

func (r *Reflect) ChainAppend(v reflect.Value) {
	r.ChainValue().Set(reflect.Append(r.ChainValue(), v))
}

func (r *Reflect) ValueIndex(i int) reflect.Value {
	return r.ChainValue().Index(i)
}

func (r *Reflect) Type() reflect.Type {
	return r.containerType().Elem()
}

func (r *Reflect) Value() reflect.Value {
	if r.IsChain() {
		log.Fatalln("model is a chain, cannot extract a value")
	}
	return r.containerValue().Elem()
}

// || STRUCT METHODS ||

func (r *Reflect) StructFieldByName(name string) reflect.Value {
	return r.Value().FieldByName(name)
}

func (r *Reflect) StructFieldByIndex(i int) reflect.Value {
	return r.Value().Field(i)

}

func (r *Reflect) StructNumFields() int {
	return r.Value().NumField()
}

// || CONSTRUCTOR ||

func (r *Reflect) NewModel() *Reflect {
	if r.IsChain() {
		return NewReflect(reflect.New(r.Type().Elem().Elem()).Interface())
	}
	return NewReflect(r.Type())
}

func (r *Reflect) NewChain() *Reflect {
	return NewReflect(reflect.MakeSlice(reflect.SliceOf(r.Type()), 0, 0))
}

func validateSliceOrStruct(v interface{}) error {
	r := v.(*Reflect)
	if !r.IsStruct() && !r.IsStruct() {
		return NewError(ErrTypeNonStructOrSlice)
	}
	return nil
}

func validateContainerIsPointer(v interface{}) error {
	r := v.(*Reflect)
	if r.containerType().Kind() != reflect.Ptr {
		return NewError(ErrTypeNonPointer)
	}
	return nil
}

var validator = validate.New([]validate.ValidateFunc{
	validateContainerIsPointer,
	validateSliceOrStruct,
})
