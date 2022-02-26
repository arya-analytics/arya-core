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

// ||| PKC |||

// PK wraps the primary key of a model,
// and provides a variety of utilities for manipulating it.
// pkChain are best created in one of two ways. Directly
// instantiate by calling:
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

// NewFromString creates a new PK with the same type as PK from the provided string.
// Returns an error fi the PK can't be converted.
func (pk PK) NewFromString(pkStr string) (PK, error) {
	var (
		newRawPk interface{}
		err      error
	)
	switch t := pk.raw.(type) {
	case uuid.UUID:
		newRawPk, err = uuid.Parse(pkStr)
	case int:
		newRawPk, err = strconv.Atoi(pkStr)
	case int32:
		newRawPk, err = strconv.Atoi(pkStr)
		newRawPk = int32(newRawPk.(int))
	case int64:
		newRawPk, err = strconv.Atoi(pkStr)
		newRawPk = int64(newRawPk.(int))
	case string:
		newRawPk = pkStr
	default:
		panic(fmt.Sprintf("pk could not be converted from string to %t. pkStr is %s", t, pkStr))
	}
	return NewPK(newRawPk), err
}

// Raw returns raw value of the pk i.e. the pk provided when calling model.NewPK.
func (pk PK) Raw() interface{} {
	return pk.raw
}

// Equals compares a provided PK tPk with PK itself. Returns true if they are equal,
// returns false if they aren't.
func (pk PK) Equals(tPk PK) bool {
	return pk.String() == tPk.String()
}

// Value returns the reflect.Value of the PK.
func (pk PK) Value() reflect.Value {
	return reflect.ValueOf(pk.raw)
}

// Type returns the type of the PK.
func (pk PK) Type() reflect.Type {
	return reflect.TypeOf(pk.raw)
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

// Strings returns the PK.String values of all PK in the PKChain.
func (pkc PKChain) Strings() (pks []string) {
	for _, pk := range pkc {
		pks = append(pks, pk.String())
	}
	return pks
}

// AllZero checks if all the PK in PKChain are the zero value.
func (pkc PKChain) AllZero() bool {
	allZero := true
	for _, pk := range pkc {
		if !pk.IsZero() {
			allZero = false
		}
	}
	return allZero
}
