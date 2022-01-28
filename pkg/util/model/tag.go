package model

import (
	"reflect"
	"strings"
)

type StructTags []StructTag

func (s StructTags) Retrieve(cat string, key string, value string) (StructTag, bool) {
	for _, st := range s {
		if st.Match(cat, key, value) {
			return st, true
		}
	}
	return StructTag{}, false
}

type StructTag struct {
	reflect.StructTag
	FldName string
}

const (
	kvSeparator = ":"
)

func constructKVPair(key string, value string) string {
	return key + kvSeparator + value
}

func (s StructTag) kvs(cat string) (kvs []string, ok bool) {
	valString, ok := s.Lookup(cat)
	if !ok {
		return kvs, ok
	}
	kvs = strings.Split(valString, ",")
	return kvs, true
}

func (s StructTag) Match(cat string, key string, value string) bool {
	kvs, ok := s.kvs(cat)
	if !ok {
		return false
	}
	for _, kv := range kvs {
		if key != "" && value != "" {
			cKv := constructKVPair(key, value)
			return cKv == kv
		} else if key != "" {
			return strings.Contains(kv, key)
		} else if value != "" {
			return strings.Contains(kv, value)
		}
		// If both our key and value are empty, then we're just looking
		// up by category, so we return true.
		return true
	}
	return false
}

func NewTags(t reflect.Type) (tags StructTags) {
	if t.Kind() != reflect.Struct {
		panic("received non-struct")
	}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tags = append(tags, StructTag{StructTag: fld.Tag, FldName: fld.Name})
	}
	return tags
}
