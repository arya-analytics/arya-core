package model

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
	"strings"
)

const (
	TagCat  = "model"
	RoleKey = "role"
	PKRole  = "pk"
)

// Reflect wraps a model object and provides utilities for accessing and manipulating
// its values. A model object is either a pointer to a struct or pointer to a slice of
// structs.  Reflect is optimal for use cases involving working with arbitrary struct
// types. It shouldn't be used in cases where only one struct type is involved.
// Instantiate Reflect by calling NewReflect.
// Avoid calling UnsafeNewReflect or instantiating Reflect directly,
// as this bypasses validation checks we execute for runtime security purposes.
type Reflect struct {
	modelObj interface{}
}

// NewReflect initializes, validates and returns a new model Reflect.
// Expects a pointer to a struct or a pointer to a slice of structs.
// Will panic if it does not receive these.
func NewReflect(modelPtr interface{}) *Reflect {
	r := UnsafeNewReflect(modelPtr)
	r.Validate()
	return r
}

// UnsafeNewReflect initializes and returns an unvalidated model
// Reflect. If you don't have a good reason to do this, don't.
// The main reason for bypassing validation is to construct a pointer from a
// provided value - see Reflect.ToNewPointer.
func UnsafeNewReflect(modelPtr interface{}) *Reflect {
	return &Reflect{modelObj: modelPtr}
}

// Validate runs validation checks against the Reflect.
// Checks that the model object is either a model or a chain.
// Panics if it isn't either of those.
func (r *Reflect) Validate() {
	if err := validator.Exec(r); err != nil {
		panic(err)
	}
}

// IsPointer returns true if the model object is a pointer to an
// object.
func (r *Reflect) IsPointer() bool {
	return r.PointerType().Kind() == reflect.Ptr
}

// Pointer returns the pointer to the Reflect model object.
func (r *Reflect) Pointer() interface{} {
	return r.modelObj
}

// || TYPE CHECKING ||

// Type returns the model object's type i.e. the actual type of the model struct.
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

// IsChain Returns true if the model object's type is chain.
func (r *Reflect) IsChain() bool {
	return r.RawType().Kind() == reflect.Slice
}

// IsStruct Returns true if the model object's type is a single struct.
func (r *Reflect) IsStruct() bool {
	return r.RawType().Kind() == reflect.Struct
}

// || STRUCT METHODS ||

// StructValue returns the value of the model object.
// Panics if the model object is a chain.
//
// This operation would panic:
// 		rChain := model.NewReflect(&[]*ExampleModel{})
// 		rChain.StructValue()
func (r *Reflect) StructValue() reflect.Value {
	r.panicIfChain()
	return r.PointerValue().Elem()
}

// StructFieldByRole retrieves the field from the model object by its role.
// Panics if the field can't be found.
func (r *Reflect) StructFieldByRole(role string) reflect.Value {
	tag, ok := r.StructTagChain().Retrieve(TagCat, RoleKey, role)
	if !ok {
		panic(fmt.Sprintf("could not find field with role %s", role))
	}
	return r.StructValue().FieldByIndex(tag.Field.Index)
}

func (r *Reflect) StructFieldByName(name string) reflect.Value {
	splitNames := strings.Split(name, ".")
	var fld = r.StructValue().FieldByName(splitNames[0])
	for _, splitName := range splitNames[1:] {
		fld = fld.Elem().FieldByName(splitName)
	}
	return fld
}

// || CHAIN METHODS ||

// ChainValue returns the value of the value of the model object.
// Panics if the model object is a struct.
//
// This operation would panic:
// 		rStruct := model.UnsafeNewReflect(&ExampleStruct{})
//		rStruct.ChainValue()
func (r *Reflect) ChainValue() reflect.Value {
	r.panicIfStruct()
	return r.RawValue()
}

// ChainAppend appends another Reflect to the model object.
// Panics if the model object is a struct.
//
// Panics if Reflect to append rta is a chain.
//
// This operation would panic:
// 		rStruct := model.UnsafeNewReflect(&ExampleStruct{})
//		rStructToAdd := model.UnsafeNewReflect(&ExampleStruct{})
//		rStruct.ChainAppend(rStructToAdd)
func (r *Reflect) ChainAppend(rta *Reflect) {
	rta.panicIfChain()
	r.panicIfStruct()
	r.ChainValue().Set(reflect.Append(r.ChainValue(), rta.PointerValue()))
}

// ChainValueByIndex retrieves Reflect from the model objet by index.
// Panics if the model object is a struct.
func (r *Reflect) ChainValueByIndex(i int) *Reflect {
	return NewReflect(r.ChainValue().Index(i).Interface())
}

// ChainValueByIndexOrNew retrieves Reflect from the model object by index.
// If the index requested exceeds the length of the chain value,
// creates new Reflect and appends it to chain value before returning.
func (r *Reflect) ChainValueByIndexOrNew(i int) *Reflect {
	diff := i - r.ChainValue().Len()
	if diff < 0 {
		return r.ChainValueByIndex(i)
	} else {
		rfl := r.NewStruct()
		r.ChainAppend(rfl)
		return rfl
	}
}

func (r *Reflect) Fields(i int) *Fields {
	var rawFlds []reflect.Value
	r.ForEach(func(rfl *Reflect, i int) {
		rawFlds = append(rawFlds, rfl.StructValue().Field(i))
	})
	t := r.Type().Field(i)
	return &Fields{
		t:      t.Type,
		values: rawFlds,
	}
}

func (r *Reflect) FieldsByName(name string) *Fields {
	var rawFlds []reflect.Value
	r.ForEach(func(rfl *Reflect, i int) {
		rawFlds = append(rawFlds, rfl.StructFieldByName(name))
	})
	t, ok := r.Type().FieldByName(name)
	if !ok {
		panic("field does not exist!")
	}
	return &Fields{
		t:      t.Type,
		values: rawFlds,
	}
}

// || FINDING VALUES ||

// ValueByPK retrieves Reflect by its pk value. If Reflect model object is chain,
// searches the chain for the PK, and returns ok=false if it can't be found.
// If Reflect model object is struct, returns the struct if the PK matches. If it isn't,
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

const forEachIfStructIndex = -1

// ForEach iterates through the model object in Reflect and calls the provided
// function. The function receives the model Reflect as well as its index.
// NOTE: The index provided to the ForEachFunc is -1 if the Reflect model object is struct.
func (r *Reflect) ForEach(fef func(rfl *Reflect, i int)) {
	if r.IsStruct() {
		fef(r, forEachIfStructIndex)
	} else {
		for i := 0; i < r.ChainValue().Len(); i++ {
			rfl := r.ChainValueByIndex(i)
			fef(rfl, i)
		}
	}
}

// || CONSTRUCTOR ||

// NewRaw creates new Reflect with the same RawType as the model object.
func (r *Reflect) NewRaw() *Reflect {
	if r.IsChain() {
		return r.NewChain()
	}
	return r.NewStruct()
}

// NewStruct creates new Reflect with a struct model object.
func (r *Reflect) NewStruct() *Reflect {
	return NewReflect(reflect.New(r.Type()).Interface())
}

// NewChain creates new Reflect with an chain model object.
func (r *Reflect) NewChain() *Reflect {
	ns := reflect.MakeSlice(reflect.SliceOf(r.NewStruct().PointerType()), 0, 0)
	p := reflect.New(ns.Type())
	p.Elem().Set(ns)
	return NewReflect(p.Interface())
}

// ToNewPointer takes the Reflect model object, creates a pointer to it,
// and creates new Reflect to the created pointer.
// Very useful for turning a struct or slice into a pointer to a struct or slice.
// It's important to validate that the reflection is not a pointer before calling this
// method, as to avoid creating pointers to pointers.
// Call this method with caution.
func (r *Reflect) ToNewPointer() *Reflect {
	p := reflect.New(r.PointerType())
	p.Elem().Set(r.PointerValue())
	return NewReflect(p.Interface())
}

// || PK ||

// PKField returns the primary key field of the model object (ie assigned role:pk).
// Panics if the field does not exist, or if the Reflect model object is a chain.
func (r *Reflect) PKField() reflect.Value {
	return r.StructFieldByRole(PKRole)
}

// PK returns new PK representing the primary key of the model.
// Panics if the field does not exist, or if the Reflect model object is a chain.
func (r *Reflect) PK() PK {
	return NewPK(r.PKField().Interface())
}

// PKChain returns all PKS in the Reflect model object.
// If the Reflect model object is a chain, returns all PK of the models in the chain.
// If Reflect model object is a struct, returns a PKChain ith length 1 containing the
// structs PK.
func (r *Reflect) PKChain() PKChain {
	var pks PKChain
	r.ForEach(func(rfl *Reflect, i int) {
		pks = append(pks, rfl.PK())
	})
	return pks
}

// || TYPE + VALUE ACCESSORS ||

// PointerType returns the type of the pointer to the model object.
// NOTE: the reflect.Kind might not be reflect.Ptr
// if a pointer wasn't provided when calling UnsafeNewReflect.
func (r *Reflect) PointerType() reflect.Type {
	return reflect.TypeOf(r.modelObj)
}

// PointerValue returns the reflect.Value of the pointer to the model object.
func (r *Reflect) PointerValue() reflect.Value {
	return reflect.ValueOf(r.modelObj)
}

// RawType returns the unparsed type of the model object.
func (r *Reflect) RawType() reflect.Type {
	return r.PointerType().Elem()
}

// RawValue returns the unparsed value of the model object.
func (r *Reflect) RawValue() reflect.Value {
	return r.PointerValue().Elem()
}

// || TAGS ||

// StructTagChain returns a StructTagChain representing all struct tags
// on the model objects type (i.e. Reflect.Type).
func (r *Reflect) StructTagChain() StructTagChain {
	return NewStructTagChain(r.Type())
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

// |||| VALIDATION ||||

func validateSliceOrStruct(v interface{}) error {
	r := v.(*Reflect)
	if !r.IsStruct() && !r.IsChain() {
		return fmt.Errorf("model validation failed - is %s, must be struct or slice",
			r.Type().Kind())
	}
	return nil
}

func validateIsPointer(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerType().Kind() != reflect.Ptr {
		return fmt.Errorf("model validation failed - model is not a pointer")
	}
	return nil
}

func validateNonZero(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerValue().IsZero() {
		return fmt.Errorf("model validation failed - model is nil")
	}
	return nil
}

var validator = validate.New([]validate.Func{
	validateIsPointer,
	validateSliceOrStruct,
	validateNonZero,
})
