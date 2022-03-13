package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Catalog []interface{}

func (mc Catalog) New(sourcePtr interface{}) interface{} {
	sourceRfl := NewReflect(sourcePtr)
	newRfl, ok := mc.retrieveCM(sourceRfl.Type())
	if !ok {
		panic(fmt.Sprintf("model %s could not be found in catalog", sourceRfl.Type().Name()))
	}
	if sourceRfl.IsChain() {
		return newRfl.NewChain().Pointer()
	}
	return newRfl.NewStruct().Pointer()
}

func (mc Catalog) Contains(sourcePtr interface{}) bool {
	_, ok := mc.retrieveCM(NewReflect(sourcePtr).Type())
	return ok
}

func (mc Catalog) retrieveCM(t reflect.Type) (*Reflect, bool) {
	for _, opt := range mc {
		optRfl := NewReflect(opt)
		if strings.EqualFold(t.Name(), optRfl.Type().Name()) {
			return optRfl, true
		}
	}
	return nil, false
}

type DataSource map[reflect.Type]*Reflect

func (ds DataSource) Retrieve(t reflect.Type) *Reflect {
	_, ok := ds[t]
	if !ok {
		ds[t] = NewReflectFromType(t).NewChain()
	}
	return ds[t]
}
