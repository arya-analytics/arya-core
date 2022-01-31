package model

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
)

// ||| PK |||

type PKChain []PK

func (pkc PKChain) Interface() (pks []interface{}) {
	for _, pk := range pkc {
		pks = append(pks, pk.Interface())
	}
	return pks
}

type PK struct {
	base interface{}
}

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

func NewPK(pk interface{}) PK {
	return PK{base: pk}
}

// String stringifies and returns the pk
func (pk PK) String() string {
	switch pk.base.(type) {
	case uuid.UUID:
		return pk.base.(uuid.UUID).String()
	case int:
		return strconv.Itoa(pk.base.(int))
	case int32:
		return strconv.Itoa(int(pk.base.(int32)))
	case int64:
		return strconv.Itoa(int(pk.base.(int64)))
	case string:
		return pk.base.(string)
	}
	panic(fmt.Sprintf("pk has unknown type %T, could not stringify", pk.base))
}

func (pk PK) Interface() interface{} {
	return pk.base
}

func (pk PK) Equals(tPk PK) bool {
	return pk.base == tPk.base
}

func (pk PK) Value() reflect.Value {
	return reflect.ValueOf(pk.base)
}

func (pk PK) IsZero() bool {
	return pk.Value().IsZero()
}
