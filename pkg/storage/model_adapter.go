package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
)

/// |||| CATALOG ||||

type ModelCatalog []reflect.Type

func (mc ModelCatalog) New(m interface{}) interface{} {
	mn := modelT(m).Name()
	for _, cm := range mc {
		if cm.Name() == mn {
			return reflect.New(cm).Interface()
		}
	}
	log.Fatalf("model %s could not be found in catalog. This is an no-op.", mn)
	return nil
}

// |||| BASE ADAPTER ||||

type ModelAdapter interface {
	Source() interface{}
	Dest() interface{}
	ExchangeToSource() error
	ExchangeToDest() error
}

func NewModelAdapter(source interface{}, dest interface{}) (ModelAdapter, error) {
	err := validateModel(source)
	if err != nil {
		return nil, err
	}
	err = validateModel(dest)
	if err != nil {
		return nil, err
	}
	sMtk, dMtk := modelT(source).Kind(), modelT(dest).Kind()
	if sMtk != dMtk {
		return nil, fmt.Errorf("models must be of the same type. Received %s and %s",
			sMtk, dMtk)
	}
	if dMtk == reflect.Struct {
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
			reflect.ValueOf(sm.Dest())))
	}
	return nil
}

func (ma *chainModelAdapter) ExchangeToSource() error {
	return ma.exchange(true)
}

func (ma *chainModelAdapter) ExchangeToDest() error {
	return ma.exchange(false)
}

func (ma *chainModelAdapter) Source() interface{} {
	return reflect.ValueOf(ma.source).Elem().Interface()
}

func (ma *chainModelAdapter) Dest() interface{} {
	return reflect.ValueOf(ma.dest).Elem().Interface()
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

func (ma *singleModelAdapter) Source() interface{} {
	return ma.sourceAm.model
}

func (ma *singleModelAdapter) Dest() interface{} {
	return ma.destAm.model
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
				return fmt.Errorf("(%s) invalid type %v for field '%s' with type %v "+
					"this is a no-op", dv.Type(), vt, k, ft)

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
	v = reflect.ValueOf(ma.Dest())
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
	v = reflect.ValueOf(ma.Dest())
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

// || VALIDATION ||
func validateModel(m interface{}) error {
	ctk := containerT(m).Kind()
	if ctk != reflect.Ptr {
		return fmt.Errorf("model container must be a pointer. received kind %s",
			containerT(m).Kind())
	}
	mtk := modelT(m).Kind()
	if mtk != reflect.Struct && mtk != reflect.Slice {
		return fmt.Errorf("model must be a struct or slice. received kind %s", mtk)
	}
	mtv := modelV(m)
	if !mtv.CanSet() {
		return fmt.Errorf("cannot set attributes on model %s", mtv)
	}
	return nil
}
