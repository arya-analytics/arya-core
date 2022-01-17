package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	log "github.com/sirupsen/logrus"
	"reflect"
)

var modelValidator = validate.New([]validate.ValidateFunc{
	validateContainerIsPointer,
	validateSliceOrStruct,
})

// TODO: Move general functions to utilities
// TODO: Figure out how to simplify type getting system
// TODO: Document APIs

/// |||| CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(m interface{}) interface{} {
	mt := modelT(m)
	mn := mt.Name()
	c := IsChainModel(mt)
	if c {
		mn = mt.Elem().Elem().Name()
	}
	for _, cm := range mc {
		tcm := reflect.TypeOf(cm)
		if tcm.Name() == mn {
			rv := reflect.New(tcm)
			if c {
				newSlice := reflect.MakeSlice(reflect.SliceOf(rv.Type()), 0, 0)
				r := reflect.New(newSlice.Type())
				r.Elem().Set(newSlice)
				return r.Interface()
			}
			return rv.Interface()
		}
	}
	log.Fatalf("model %s could not be found in catalog. This is an no-op.", mn)
	return nil
}

// |||| BASE ADAPTER ||||

type ModelAdapter interface {
	SourcePointer() interface{}
	DestPointer() interface{}
	DestValue() reflect.Value
	SourceValue() reflect.Value
	ExchangeToSource() error
	ExchangeToDest() error
}

func NewModelAdapter(source interface{}, dest interface{}) (ModelAdapter, error) {
	if err := modelValidator.Exec(source); err != nil {
		return nil, err
	}
	if err := modelValidator.Exec(dest); err != nil {
		return nil, err
	}
	sMtk, dMtk := modelT(source).Kind(), modelT(dest).Kind()
	if sMtk != dMtk {
		return nil, NewError(ErrTypeIncompatibleModels)
	}
	if !IsChainModel(modelT(source)) {
		return newSingleModelAdapter(source, dest), nil
	}
	return &chainModelAdapter{source, dest}, nil
}

type modelValues map[string]interface{}

// |||| MULTI MODEL ADAPTER ||||

type chainModelAdapter struct {
	source interface{}
	dest   interface{}
}

func (ma *chainModelAdapter) exchange(toSource bool) error {
	var to interface{}
	var from interface{}
	if toSource {
		to, from = ma.source, ma.dest
	} else {
		to, from = ma.dest, ma.source
	}
	fromRv := modelV(from)
	toModelSliceValue := modelV(to)
	toModelType := modelT(to).Elem().Elem()
	for i := 0; i < fromRv.Len(); i++ {
		fromMv := fromRv.Index(i).Interface()
		var toMv interface{}
		if i >= toModelSliceValue.Len() {
			toMv = reflect.New(toModelType).Interface()
		} else {
			toMv = toModelSliceValue.Index(i).Interface()
		}
		sm := newSingleModelAdapter(fromMv, toMv)
		if err := sm.ExchangeToDest(); err != nil {
			return err
		}
		toModelSliceValue.Set(reflect.Append(toModelSliceValue,
			reflect.ValueOf(sm.DestPointer())))
	}
	return nil
}

func (ma *chainModelAdapter) ExchangeToSource() error {
	return ma.exchange(true)
}

func (ma *chainModelAdapter) ExchangeToDest() error {
	return ma.exchange(false)
}

func (ma *chainModelAdapter) SourcePointer() interface{} {
	return reflect.ValueOf(ma.source).Interface()
}

func (ma *chainModelAdapter) DestPointer() interface{} {
	return reflect.ValueOf(ma.dest).Interface()
}

func (ma *chainModelAdapter) DestValue() reflect.Value {
	return reflect.ValueOf(ma.dest).Elem()
}

func (ma *chainModelAdapter) SourceValue() reflect.Value {
	return reflect.ValueOf(ma.source).Elem()
}

// |||| MODEL ADAPTER ||||

type singleModelAdapter struct {
	sourceAm *adaptedModel
	destAm   *adaptedModel
}

func newSingleModelAdapter(source interface{}, dest interface{}) *singleModelAdapter {
	return &singleModelAdapter{
		sourceAm: &adaptedModel{model: source},
		destAm:   &adaptedModel{model: dest},
	}
}

func (ma *singleModelAdapter) SourcePointer() interface{} {
	return ma.sourceAm.model
}

func (ma *singleModelAdapter) DestPointer() interface{} {
	return ma.destAm.model
}

func (ma *singleModelAdapter) DestValue() reflect.Value {
	return reflect.ValueOf(ma.destAm.model)
}

func (ma *singleModelAdapter) SourceValue() reflect.Value {
	return reflect.ValueOf(ma.sourceAm.model)
}

func (ma *singleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.bindVals(ma.destAm.mapVals())
}

func (ma *singleModelAdapter) ExchangeToDest() error {
	return ma.destAm.bindVals(ma.sourceAm.mapVals())
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	model interface{}
}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) error {
	dv := modelV(mw.model)
	for k, rv := range mv {
		field := dv.FieldByName(k)
		val := reflect.ValueOf(rv)
		if !field.CanSet() {
			continue
		}
		vt, ft := val.Type(), field.Type()
		if vt != ft {
			invalid := true
			if val.Kind() == reflect.Ptr {
				invalid = false
				valModelKind := modelT(val.Interface()).Kind()
				var err error
				if valModelKind == reflect.Slice {
					if modelV(val.Interface()).Len() == 0 {
						continue
					}
					val, err = adaptSlice(val, field)
				} else if valModelKind == reflect.Struct {
					valModelVal := modelV(val.Interface())
					if !valModelVal.IsValid() {
						continue
					}
					val, err = adaptStruct(val, field)
				} else {
					invalid = true
				}
				if err != nil {
					return err
				}
			}
			if invalid {
				return NewError(ErrTypeInvalidField)
			}
		}
		field.Set(val)
	}
	return nil
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	var mv = modelValues{}
	dv := modelV(mw.model)
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		// Need to convert slices to addy because that's what NewModelAdapter expects.
		if f.Type().Kind() == reflect.Slice {
			mv[t.Name] = f.Addr().Interface()
		} else {
			mv[t.Name] = f.Interface()
		}
	}
	return mv
}

// |||| UTILITIES ||||

// || ADAPTATION ||
func adaptSlice(slc reflect.Value, fld reflect.Value) (v reflect.Value, e error) {
	ma, err := NewModelAdapter(slc.Interface(), fld.Addr().Interface())
	err = ma.ExchangeToDest()
	if err != nil {
		return v, err
	}
	v = reflect.ValueOf(ma.DestValue().Interface())
	return v, e
}

func adaptStruct(sct reflect.Value, fld reflect.Value) (v reflect.Value, e error) {
	source := sct.Interface()
	dest := reflect.New(fld.Type().Elem()).Interface()
	ma, err := NewModelAdapter(source, dest)
	err = ma.ExchangeToDest()
	if err != nil {
		return v, err
	}
	v = ma.DestValue()
	return v, e
}

// || TYPE AND VALUE ACCESS ||
func containerT(m interface{}) reflect.Type {
	return reflect.TypeOf(m)
}

func containerV(m interface{}) reflect.Value {
	return reflect.ValueOf(m)
}

func modelV(m interface{}) reflect.Value {
	return containerV(m).Elem()
}

func modelT(m interface{}) reflect.Type {
	return containerT(m).Elem()
}

func IsChainModel(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}
