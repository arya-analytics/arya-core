package storage

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

const (
	// destValueDepth is the amount of pointer dereferences needed to access the
	// underlying struct
	destValueDepth = 2
)

type ChannelConfig struct {
	ID   uuid.UUID
	Name string
}

type ModelValues map[string]interface{}

type Model struct {
	Dest interface{}
}

// BindVals binds a set of ModelValues to the Model fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (m *Model) BindVals(mv ModelValues) error {
	dv := m.destVal()
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

// MapVals maps Model fields to ModelValues.
func (m *Model) MapVals() ModelValues {
	var mv = ModelValues{}
	dv := m.destVal()
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}

// Interface returns the model field as an interface that can be rebound to its
// original struct.
func (m *Model) Interface() interface{} {
	return m.Dest
}

func (m *Model) destVal() (v reflect.Value) {
	v = reflect.ValueOf(&m.Dest)
	for i := 0; i <= destValueDepth; i++ {
		v = v.Elem()
	}
	return v
}
