package storage

import (
	"fmt"
	"reflect"
)

const (
	// destValueDepth is the amount of pointer dereferences needed to access the
	// underlying struct
	destValueDepth = 2
)

type ChannelConfig struct {
	ID   int32
	Name string
}

type ModelValues map[string]interface{}

type ModelWrapper struct {
	model interface{}
}

// NewModelWrapper creates a new ModelWrapper from a provided model struct.
func NewModelWrapper(dest interface{}) *ModelWrapper {
	return &ModelWrapper{model: dest}
}

// BindVals binds a set of ModelValues to the ModelWrapper fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *ModelWrapper) BindVals(mv ModelValues) error {
	dv := mw.destVal()
	for k, rv := range mv {
		f := dv.FieldByName(k)
		v := reflect.ValueOf(rv)
		if !f.IsValid() {
			return fmt.Errorf("invalid key %s while binding to %v", k, dv.Type())
		}
		if !f.CanSet() {
			return fmt.Errorf("unsettable key %s in vals while binding to %v ",
				k, dv.Type())
		}
		vt, ft := v.Type(), f.Type()
		if vt != ft {
			return fmt.Errorf("invalid type %v for field '%s' with type %v- "+
				"this is a no-op", vt, k, ft)
		}
		f.Set(v)
	}
	return nil
}

// MapVals maps ModelWrapper fields to ModelValues.
func (mw *ModelWrapper) MapVals() ModelValues {
	var mv = ModelValues{}
	dv := mw.destVal()
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}

// Model returns the wrappers model
func (m *ModelWrapper) Model() interface{} {
	return m.model
}

func (mw *ModelWrapper) destVal() (v reflect.Value) {
	v = reflect.ValueOf(&mw.model)
	for i := 0; i <= destValueDepth; i++ {
		v = v.Elem()
	}
	return v
}
