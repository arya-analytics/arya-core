package model

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
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

type DataSource struct {
	mu   sync.Mutex
	data map[reflect.Type]*Reflect
}

func NewDataSource() *DataSource {
	return &DataSource{
		data: make(map[reflect.Type]*Reflect),
	}
}

func (ds *DataSource) Retrieve(t reflect.Type) *Reflect {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.data[t]
	if !ok {
		ds.data[t] = NewReflectFromType(t).NewChain()
	}
	return ds.data[t]
}

func (ds *DataSource) Write(rfl *Reflect) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[rfl.Type()] = rfl
}
