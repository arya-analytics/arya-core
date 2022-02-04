// Package model holds utilities for manipulating models.
//
// What are models? They are arbitrary slices or structs.
// Their most notable attribute however,
// is that they are used as arguments to APIs that need to work with arbitrary struct
// and struct slice types.
//
// A prime example is the Arya Core storage layer,
// which uses models and this package to save and query arbitrary data types to/from
// storage services.
//
// The core functionality of this package involves wrapping model structs of slices
// using Reflect and NewReflect.
package model

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
)

// ||| PK |||

// PK wraps the primary key of a model,
// and provides a variety of utilities for manipulating it.
// pkChain are best created in one of two ways. The first,
// by directly instantiating by calling:
//
// 		model.NewPK(m.PKField)
//
// or by instantiation of model.Reflect subsequent call of Reflect.PK().
type PK struct {
	raw interface{}
}

// NewPK constructs and returns a new PK.
func NewPK(pk interface{}) PK {
	return PK{raw: pk}
}

// String converts and returns the pk as a string. Supported PK types are uuid.UUID,
// int, int32, int64, and string.
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
	panic(fmt.Sprintf("pk has unknown type %T, could not stringify", pk.raw))
}

// Raw returns raw value of the pk i.e. the pk provided when calling model.NewPK.
func (pk PK) Raw() interface{} {
	return pk.raw
}

// Equals compares a provided PK tPk with PK itself. Returns true if they are equal,
// returns false if they aren't.
func (pk PK) Equals(tPk PK) bool {
	return pk.raw == tPk.raw
}

// Value returns the reflect.Value of the PK.
func (pk PK) Value() reflect.Value {
	return reflect.ValueOf(pk.raw)
}

// IsZero returns true if the pk is the zero value for its type.
func (pk PK) IsZero() bool {
	return pk.Value().IsZero()
}

// PKChain provides utilities for managing a chain of PK.
type PKChain []PK

// NewPKChain creates a new set of primary keys from a slice of arbitrary values.
// Will panic if pks is not a slice.
func NewPKChain(pks interface{}) (pkc PKChain) {
	val := reflect.ValueOf(pks)
	if val.Type().Kind() != reflect.Slice {
		panic("model.NewPKS received non slice type")
	}
	for i := 0; i < val.Len(); i++ {
		pkc = append(pkc, NewPK(val.Index(i).Interface()))
	}
	return pkc
}

// Raw returns the PK.Raw values of all PK in PKChain.
func (pkc PKChain) Raw() (pks []interface{}) {
	for _, pk := range pkc {
		pks = append(pks, pk.Raw())
	}
	return pks
}
