package model

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
)

const (
	tagCat  = "model"
	roleKey = "role"
	pkRole  = "pk"
)

type Reflect struct {
	modelPtr interface{}
}

func NewReflect(modelPtr interface{}) *Reflect {
	return &Reflect{
		modelPtr: modelPtr,
	}
}

func (r *Reflect) Validate() {
	if err := validator.Exec(r); err != nil {
		panic(err)
	}
}

// IsPointer returns true if the value provided to Reflect is a pointer to an
// underlying model.
func (r *Reflect) IsPointer() bool {
	return r.PointerType().Kind() == reflect.Ptr
}

// Pointer returns the pointer to the underlying model in Reflect.
func (r *Reflect) Pointer() interface{} {
	return r.modelPtr
}

// || TYPE CHECKING ||

// Type returns the 'model' type i.e. the actual type of the Model Struct.
func (r *Reflect) Type() reflect.Type {
	if r.IsChain() {
		/* base type is the slice
		first elem is pointer to struct
		second elem is struct */
		return r.RawType().Elem().Elem()
	}
	/* base type is pointer to struct
	first elem is struct */
	return r.RawType()
}

// IsChain Returns true if the reflection contains a chain of model structs.
func (r *Reflect) IsChain() bool {
	return r.RawType().Kind() == reflect.Slice
}

// IsStruct Returns true if the reflection contains a single model struct.
func (r *Reflect) IsStruct() bool {
	return r.RawType().Kind() == reflect.Struct
}

// || STRUCT METHODS ||

// StructValue returns the value of the Reflect model struct.
// Panics if the reflection contains a chain.
//
// This operation would panic:
// 		rChain := model.NewReflect(&[]*ExampleModel{})
// 		rChain.StructValue()
func (r *Reflect) StructValue() reflect.Value {
	r.panicIfChain()
	return r.PointerValue().Elem()
}

// StructFieldByRole retrieves the field from the Reflect model struct by its role.
// Panics if the field can't be found.
func (r *Reflect) StructFieldByRole(role string) reflect.Value {
	tag, ok := r.StructTagChain().Retrieve(tagCat, roleKey, role)
	if !ok {
		panic(fmt.Sprintf("could not find field with role %s", role))
	}
	return r.StructValue().FieldByIndex(tag.Field.Index)
}

// || CHAIN METHODS ||

// ChainValue returns the value of the Reflect model chain.
// Panics if the reflection contains a struct.
//
// This operation would panic:
// 		rStruct := model.NewReflect(&ExampleStruct{})
//		rStruct.ChainValue()
func (r *Reflect) ChainValue() reflect.Value {
	r.panicIfStruct()
	return r.RawValue()
}

// ChainAppend appends another Reflect to the model chain.
// Panics if the reflection is a struct.
//
// Provided Reflect v must contain a struct.
//
// This operation would panic:
// 		rStruct := model.NewReflect(&ExampleStruct{})
//		rStructToAdd := model.NewReflect(&ExampleStruct{})
//		rStruct.ChainAppend(rStructToAdd)
func (r *Reflect) ChainAppend(v *Reflect) {
	r.panicIfStruct()
	r.ChainValue().Set(reflect.Append(r.ChainValue(), v.PointerValue()))
}

// ChainValueByIndex retrieves Reflect from the model chain by index.
// Panics if the reflection is a struct.
func (r *Reflect) ChainValueByIndex(i int) *Reflect {
	return NewReflect(r.ChainValue().Index(i).Interface())
}

// || FINDING VALUES ||

// ValueByPK retrieves Reflect by its pk value. If Reflect contains a chain, searches
// the chain for the PK, and returns ok=false if it can't be found.
// If Reflect contains a struct, returns the struct if the PK matches. If not,
// returns ok=false.
func (r *Reflect) ValueByPK(pk PK) (retRfl *Reflect, ok bool) {
	r.ForEach(func(rfl *Reflect, i int) {
		if rfl.PK().Equals(pk) {
			retRfl = rfl
		}
	})
	if retRfl == nil {
		return retRfl, false
	}
	return retRfl, true
}

// || ITERATION UTILITIES ||

const structIndex = -1

// ForEach iterates through each model struct in Reflect and calls the provided
// function. The function receives the model Reflect as well as its index.
// NOTE: The index provided to the ForEachFunc is -1 if the Reflect contains a struct
// internally.
func (r *Reflect) ForEach(fef func(rfl *Reflect, i int)) {
	if r.IsStruct() {
		fef(r, structIndex)
	} else {
		for i := 0; i < r.ChainValue().Len(); i++ {
			rfl := r.ChainValueByIndex(i)
			fef(rfl, i)
		}
	}
}

// || CONSTRUCTOR ||

// NewRaw creates new Reflect with the same RawType as source Reflect.
func (r *Reflect) NewRaw() *Reflect {
	if r.IsChain() {
		return r.NewChain()
	}
	return r.NewStruct()
}

// NewStruct creates new Reflect with an internal struct.
func (r *Reflect) NewStruct() *Reflect {
	return NewReflect(reflect.New(r.Type()).Interface())
}

// NewChain creates new Reflect with an internal chain.
func (r *Reflect) NewChain() *Reflect {
	ns := reflect.MakeSlice(reflect.SliceOf(r.NewStruct().PointerType()), 0, 0)
	p := reflect.New(ns.Type())
	p.Elem().Set(ns)
	return NewReflect(p.Interface())
}

// ToNewPointer takes the Reflect, creates a pointer to it,
// and creates new Reflect to the pointer.
// Very useful for turning a struct or chain into a pointer to a struct or chain.
// It's important to validate that the reflection is not a pointer before calling this
// method, as to avoid creating pointers to pointers.
func (r *Reflect) ToNewPointer() *Reflect {
	p := reflect.New(r.PointerType())
	p.Elem().Set(r.PointerValue())
	return NewReflect(p.Interface())
}

// || PK ||

// PKField returns the primary key field of the model (ie assigned role:pk).
// Panics if the field does not exist, or if the Reflect is a struct.
func (r *Reflect) PKField() reflect.Value {
	return r.StructFieldByRole(pkRole)
}

// PK returns new PK representing the primary key of the model.
// Panics if the field does not exist, or if the Reflect is a struct.
func (r *Reflect) PK() PK {
	return NewPK(r.PKField().Interface())
}

// PKChain returns all PKS in the Reflect. If the Reflect contains a chain,
// returns all PKChain of the models in the chain. If Reflect contains a struct,
// returns a slice with length 1 containing the structs PK.
func (r *Reflect) PKChain() PKChain {
	var pks PKChain
	r.ForEach(func(rfl *Reflect, i int) {
		pks = append(pks, rfl.PK())
	})
	return pks
}

// || TYPE + VALUE ACCESSORS ||

// PointerType returns the reflect.Type of the pointer to the underlying model.
// NOTE: the reflect.Kind might not be reflect.Ptr
// if a pointer wasn't provided when calling NewReflect.
func (r *Reflect) PointerType() reflect.Type {
	return reflect.TypeOf(r.modelPtr)
}

// PointerValue returns the reflect.Value of the pointer to the underlying model.
func (r *Reflect) PointerValue() reflect.Value {
	return reflect.ValueOf(r.modelPtr)
}

// RawType returns the unparsed type of the model (
// or model chain) contained within the Reflect.
func (r *Reflect) RawType() reflect.Type {
	return r.PointerType().Elem()
}

// RawValue returns the unparsed value of the model (
// or model chain) contained within the Reflect.
func (r *Reflect) RawValue() reflect.Value {
	return r.PointerValue().Elem()
}

// || TYPE ASSERTIONS ||

func (r *Reflect) panicIfChain() {
	if r.IsChain() {
		panic("model is chain, cannot get struct value")
	}
}

func (r *Reflect) panicIfStruct() {
	if r.IsStruct() {
		panic("model is struct, cannot get chain value")
	}
}

// || TAGS ||

// StructTagChain returns a set of StructTagChain representing all struct tags
// on Reflect.Type.
func (r *Reflect) StructTagChain() StructTagChain {
	return NewStructTagChain(r.Type())
}

// |||| VALIDATION ||||

func validateSliceOrStruct(v interface{}) error {
	r := v.(*Reflect)
	if !r.IsStruct() && !r.IsChain() {
		return fmt.Errorf("model validation failed, is %s must be struct or slice",
			r.Type().Kind())
	}
	return nil
}

func validateIsPointer(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerType().Kind() != reflect.Ptr {
		return fmt.Errorf("model validation failed. model is not a pointer")
	}
	return nil
}

func validateNonZero(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerValue().IsZero() {
		return fmt.Errorf("model validation failed. model is nil")
	}
	return nil
}

var validator = validate.New([]validate.Func{
	validateIsPointer,
	validateSliceOrStruct,
	validateNonZero,
})
