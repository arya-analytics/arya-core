package model

import (
	"fmt"
	"reflect"
	"strings"
)

// StructTag extends the reflect.StructTag API and provides utilities for searching
// for and manipulating struct tags. Has no constructor.
type StructTag struct {
	reflect.StructTag
	Field reflect.StructField
}

// RetrieveKVChain the key Val pairs in a struct tag category as strings i.e.
// ["ferrari:Fast","beetle:s"].
func (s StructTag) RetrieveKVChain(cat string) (kvc []string, ok bool) {
	valString, ok := s.Lookup(cat)
	if !ok {
		return kvc, false
	}
	kvc = strings.Split(valString, kvChainSeparator)
	return kvc, true
}

func (s StructTag) Retrieve(cat, key string) (string, bool) {
	kvc, ok := s.RetrieveKVChain(cat)
	if !ok {
		return "", false
	}
	for _, kv := range kvc {
		kOpt, val := splitKVPair(kv)
		if kOpt == key {
			return val, true
		}
	}
	return "", false
}

const allMatchIndicator = "*"

// Match returns true if the provided arguments match the category, key,
// and/or Val of the StructTag. If you want to search by arbitrary Val,
// pass a "*" to key arg. NOTE: "*" will not search all categories,
// and "*" will also not search all values.
func (s StructTag) Match(cat string, key string, value string) bool {
	kvs, ok := s.RetrieveKVChain(cat)
	if !ok {
		return false
	}
	var matcher func(kv string) bool
	if key != allMatchIndicator && value != allMatchIndicator {
		cKv := constructKVPair(key, value)
		matcher = func(kv string) bool { return cKv == kv }
	} else if key != allMatchIndicator {
		matcher = func(kv string) bool { return strings.Contains(kv, key) }
	} else if value != allMatchIndicator {
		panic("cannot search struct tag by arbitrary Val")
	} else {
		// If both our key and Val are empty, then we're just looking
		// up by category, so we return true.
		return true
	}
	for _, kv := range kvs {
		if matcher(kv) {
			return true
		}
	}
	return false
}

// StructTagChain provides utilities for managing a chain of struct tags.
//This data type is ideal for storing the StructTags of the fields of a model.
type StructTagChain []StructTag

// NewStructTagChain creates a new chain of struct tags from the provided struct
// type. Creates and appends a new StructTag for each field in the struct.
// Panics of the provided reflect.Type t is not of reflect.Struct kind.
func NewStructTagChain(t reflect.Type) (tags StructTagChain) {
	if t.Kind() != reflect.Struct {
		panic("model.NewStructTagChain - received non-struct")
	}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tags = append(tags, StructTag{StructTag: fld.Tag, Field: fld})
	}
	return tags
}

// Retrieve retrieves the StructTag in the StructTagChain that matches the provided arguments.
// If you want to search by arbitrary Val, pass a "*" to key. If you want to search by
// arbitrary key, pass a "*" to key arg. NOTE: "*" will not search all categories.
func (s StructTagChain) Retrieve(cat string, key string, value string) (StructTag, bool) {
	for _, st := range s {
		if st.Match(cat, key, value) {
			return st, true
		}
	}
	return StructTag{}, false
}

// RetrieveByFieldName retrieves the struct tag in the StructTagChain by its field name.
func (s StructTagChain) RetrieveByFieldName(fldName string) (StructTag, bool) {
	for _, st := range s {
		if fieldNamesEqual(st.Field.Name, fldName) {
			return st, true
		}
	}
	return StructTag{}, false
}

// RetrieveByFieldRole retrieves a field by its role struct tag.
func (s StructTagChain) RetrieveByFieldRole(role string) (StructTag, bool) {
	return s.Retrieve(TagCat, RoleKey, role)
}

// RetrieveBase retrieves the model.Base tag for hte field.
const baseFieldName = "Base"

func (s StructTagChain) RetrieveBase() (StructTag, bool) {
	return s.RetrieveByFieldName(baseFieldName)
}

// HasAnyFields determines if the chain contains and struct fields with a name in the provided WhereFields.
func (s StructTagChain) HasAnyFields(flds ...string) bool {
	for _, st := range s {
		for _, fld := range flds {
			if SplitFirstFieldName(fld) == st.Field.Name {
				return true
			}
		}
	}
	return false
}

const (
	kvPairSeparator  = ":"
	kvChainSeparator = ","
)

func constructKVPair(key string, value string) string {
	return key + kvPairSeparator + value
}

func splitKVPair(kvp string) (string, string) {
	s := strings.Split(kvp, kvPairSeparator)
	if len(s) != 2 {
		panic(fmt.Sprintf("key value pair %s improperly formatted", kvp))
	}
	return s[0], s[1]
}
