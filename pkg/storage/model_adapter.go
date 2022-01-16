package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
)

const (
	// destValueDepth is the amount of pointer dereferences needed to access the
	// underlying struct
	destValueDepth = 2
)

type ModelAdapter interface {
	Source() interface{}
	Dest() interface{}
	ExchangeToSource() error
	ExchangeToDest() error
}

func NewModelAdapter(source interface{}, dest interface{}) ModelAdapter {
	sourceKind := reflect.TypeOf(source).Elem().Kind()
	destKind := reflect.TypeOf(dest).Elem().Kind()
	if sourceKind != destKind {
		log.Fatalln("Source kind is not equal to dest kind")

	}
	if sourceKind == reflect.Slice {
		return NewMultiModelAdapter(source, dest)
	}
	return NewSingleModelAdapter(source, dest)
}

type ModelValues map[string]interface{}

// |||| MULTI MODEL ADAPTER ||||

type MultiModelAdapter struct {
	sourceModels interface{}
	destModels   interface{}
}

func NewMultiModelAdapter(sourceModels interface{},
	destModels interface{}) *MultiModelAdapter {
	return &MultiModelAdapter{
		sourceModels: sourceModels,
		destModels:   destModels,
	}
}

func (ma *MultiModelAdapter) exchange(toSource bool) error {
	var exchangeTo interface{}
	var exchangeFrom interface{}
	if toSource {
		exchangeTo = ma.destModels
		exchangeFrom = ma.sourceModels
	} else {
		exchangeTo = ma.sourceModels
		exchangeFrom = ma.destModels
	}
	destRv := reflect.ValueOf(exchangeTo).Elem()
	sourceRv := reflect.ValueOf(exchangeFrom).Elem()
	if destRv.Kind() == reflect.Slice {
		for i := 0; i < destRv.Len(); i++ {
			destMv := destRv.Index(i).Interface()
			var sourceMv interface{}
			if i >= sourceRv.Len() {
				sourceMv = reflect.New(sourceRv.Type().Elem().Elem()).Interface()
			} else {
				sourceMv = sourceRv.Index(i).Interface()
			}
			ma := NewSingleModelAdapter(sourceMv, destMv)
			if err := ma.ExchangeToSource(); err != nil {
				return err
			}
			sourceRv.Set(reflect.Append(sourceRv, reflect.ValueOf(ma.Source())))
		}
	}
	return nil
}

func (ma *MultiModelAdapter) ExchangeToSource() error {
	return ma.exchange(true)

}

func (ma *MultiModelAdapter) ExchangeToDest() error {
	return ma.exchange(false)
}

func (ma *MultiModelAdapter) Source() interface{} {
	return reflect.ValueOf(ma.sourceModels).Elem().Interface()
}

func (ma *MultiModelAdapter) Dest() interface{} {
	return reflect.ValueOf(ma.Source()).Elem().Interface()
}

// |||| MODEL ADAPTER ||||

type SingleModelAdapter struct {
	sourceAm *AdaptedModel
	destAm   *AdaptedModel
}

func NewSingleModelAdapter(sourceModel interface{}, destModel interface{}) *SingleModelAdapter {
	return &SingleModelAdapter{
		sourceAm: NewAdaptedModel(sourceModel),
		destAm:   NewAdaptedModel(destModel),
	}
}

func (ma *SingleModelAdapter) Source() interface{} {
	return ma.sourceAm.Model()
}

func (ma *SingleModelAdapter) Dest() interface{} {
	return ma.destAm.model
}

func (ma *SingleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.BindVals(ma.destAm.MapVals())
}

func (ma *SingleModelAdapter) ExchangeToDest() error {
	return ma.destAm.BindVals(ma.sourceAm.MapVals())
}

// |||| ADAPTED MODEL |||||

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
			return fmt.Errorf("(%s) invalid type %v for field '%s' with type %v- "+
				"this is a no-op", dv.Type(), vt, k, ft)
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
