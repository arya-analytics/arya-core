package model

import (
	"reflect"
	"strings"
)

type StructTagChain []StructTag

func (s StructTagChain) Retrieve(cat string, key string, value string) (StructTag, bool) {
	for _, st := range s {
		if st.match(cat, key, value) {
			return st, true
		}
	}
	return StructTag{}, false
}

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

type StructTag struct {
	reflect.StructTag
	Field reflect.StructField
}

const (
	kvPairSeparator  = ":"
	kvChainSeparator = ","
)

func constructKVPair(key string, value string) string {
	return key + kvPairSeparator + value
}

func (s StructTag) retrieveKVChain(cat string) (kvc []string, ok bool) {
	valString, ok := s.Lookup(cat)
	if !ok {
		return kvc, false
	}
	kvc = strings.Split(valString, kvChainSeparator)
	return kvc, true
}

func (s StructTag) match(cat string, key string, value string) bool {
	kvs, ok := s.retrieveKVChain(cat)
	if !ok {
		return false
	}
	var matcher func(kv string) bool
	if key != "*" && value != "*" {
		cKv := constructKVPair(key, value)
		matcher = func(kv string) bool { return cKv == kv }
	} else if key != "*" {
		matcher = func(kv string) bool { return strings.Contains(kv, key) }
	} else if value != "*" {
		panic("StructTag.match - can't perform a tag match with value and no key.")
	} else {
		// If both our key and value are empty, then we're just looking
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
