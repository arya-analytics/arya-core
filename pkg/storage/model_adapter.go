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

type ModelValues map[string]interface{}

// |||| MODEL ADAPTER ||

type ModelAdapter struct {
	sourceAm *AdaptedModel
	destAm   *AdaptedModel
}

func NewModelAdapter(sourceModel interface{}, destModel interface{}) *ModelAdapter {
	return &ModelAdapter{
		sourceAm: NewAdaptedModel(sourceModel),
		destAm:   NewAdaptedModel(destModel),
	}
}

func (ma *ModelAdapter) SourceModel() interface{} {
	return ma.sourceAm.Model()
}

func (ma *ModelAdapter) DestModel() interface{} {
	return ma.destAm.model
}

func (ma *ModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.BindVals(ma.destAm.MapVals())
}

func (ma *ModelAdapter) ExchangeToDest() error {
	return ma.destAm.BindVals(ma.sourceAm.MapVals())
}

// || ADAPTED MODEL ||

type AdaptedModel struct {
	model interface{}
}

// NewAdaptedModel creates a new AdaptedModel from a provided model struct.
func NewAdaptedModel(model interface{}) *AdaptedModel {
	return &AdaptedModel{model: model}
}

// BindVals binds a set of ModelValues to the AdaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *AdaptedModel) BindVals(mv ModelValues) error {
	dv := mw.destVal()
	for k, rv := range mv {
		f := dv.FieldByName(k)
		v := reflect.ValueOf(rv)
		if !f.IsValid() {
			return fmt.Errorf("invalid key %storage while binding to %v", k, dv.Type())
		}
		if !f.CanSet() {
			return fmt.Errorf("unsettable key %storage in vals while binding to %v ",
				k, dv.Type())
		}
		vt, ft := v.Type(), f.Type()
		if vt != ft {
			return fmt.Errorf("invalid type %v for field '%storage' with type %v- "+
				"this is a no-op", vt, k, ft)
		}
		f.Set(v)
	}
	return nil
}

// MapVals maps AdaptedModel fields to ModelValues.
func (mw *AdaptedModel) MapVals() ModelValues {
	var mv = ModelValues{}
	dv := mw.destVal()
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}

// Model returns the wrapped model
func (mw *AdaptedModel) Model() interface{} {
	return mw.model
}

func (mw *AdaptedModel) destVal() (v reflect.Value) {
	v = reflect.ValueOf(&mw.model)
	for i := 0; i <= destValueDepth; i++ {
		v = v.Elem()
	}
	return v
}
