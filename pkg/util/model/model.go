package model

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/google/uuid"
	"reflect"
	"strconv"
)

const (
	TagCat  = "model"
	RoleKey = "role"
	PKRole  = "pk"
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
	err := validator.Exec(r)
	if err != nil {
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
		/* raw type is the slice
		first elem is pointer to struct
		second elem is struct */
		return r.RawType().Elem().Elem()
	}
	/* raw type is pointer to struct
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
	tag, ok := r.Tags().Retrieve(TagCat, RoleKey, role)
	if !ok {
		panic(fmt.Sprintf("Could not find field with role %s", role))
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

const ifStructIndex = -1

// ForEachFunc is called in Reflect.ForEach,
// and provides the model reflection as well as its index.
type ForEachFunc func(rfl *Reflect, i int)

// ForEach iterates through each model struct in Reflect and calls the provided
// ForEachFunc.
// NOTE: The index provided to the ForEachFunc is -1 if the Reflect contains a struct
// internally.
func (r *Reflect) ForEach(fef ForEachFunc) {
	if r.IsStruct() {
		fef(r, ifStructIndex)
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
	return r.StructFieldByRole(PKRole)
}

// PK returns new PK representing the primary key of the model.
// Panics if the field does not exist, or if the Reflect is a struct.
func (r *Reflect) PK() PK {
	return NewPK(r.PKField().Interface())
}

// PKs returns all PKS in the Reflect. If the Reflect contains a chain,
// returns all PKs of the models in the chain. If Reflect contains a struct,
// returns a slice with length 1 containing the structs PK.
func (r *Reflect) PKs() PKs {
	var pks PKs
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

// ValueForSet is useful for getting the reflect.
// Value required when setting on a model field.
func (r *Reflect) ValueForSet() reflect.Value {
	if r.IsChain() {
		return r.RawValue()
	}
	return r.PointerValue()
}

// || TYPE ASSERTIONS ||

func (r *Reflect) panicIfChain() {
	if r.IsChain() {
		panic("model is a chain, cannot get a struct value")
	}
}

func (r *Reflect) panicIfStruct() {
	if r.IsStruct() {
		panic("model is struct, cannot get a chain value")
	}
}

// || TAGS ||

// Tags returns a set of StructTags representing all struct tags on Reflect.Type.
func (r *Reflect) Tags() StructTags {
	return NewTags(r.Type())
}

// |||| VALIDATION ||||

// || REFLECT ||

func validateSliceOrStruct(v interface{}) error {
	r := v.(*Reflect)
	if !r.IsStruct() && !r.IsChain() {
		return fmt.Errorf("model reflect validation failed. " +
			"the provided model is not a pointer")
	}
	return nil
}

func validateIsPointer(v interface{}) error {
	r := v.(*Reflect)
	if r.PointerType().Kind() != reflect.Ptr {
		return fmt.Errorf("model reflect validation failed. " +
			"the provided model is not a pointer")
	}
	return nil
}

var validator = validate.New([]validate.Func{
	validateIsPointer,
	validateSliceOrStruct,
})

// ||| PK |||

type PKs []PK

func (pks PKs) Interface() (pkis []interface{}) {
	for _, pk := range pks {
		pkis = append(pkis, pk.Interface())
	}
	return pkis
}

type PK struct {
	raw interface{}
}

func NewPK(pk interface{}) PK {
	return PK{raw: pk}
}

func (pk PK) String() string {
	switch pk.raw.(type) {
	case uuid.UUID:
		return pk.raw.(uuid.UUID).String()
	case int:
		return strconv.Itoa(pk.raw.(int))
	case int32:
		return strconv.Itoa(int(pk.raw.(int32)))
	case int64:
		return strconv.Itoa(int(pk.raw.(int64)))
	case string:
		return pk.raw.(string)
	}
	panic("Could not convert PK to string")
}

func (pk PK) Interface() interface{} {
	return pk.raw
}

func (pk PK) Equals(tPk PK) bool {
	return pk.raw == tPk.raw
}

func (pk PK) Value() reflect.Value {
	return reflect.ValueOf(pk.raw)
}

func (pk PK) IsZero() bool {
	return pk.Value().IsZero()
}
