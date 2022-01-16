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

// |||| ADAPTER OPTS ||||

type ModelAdapterOpts struct {
	Source        interface{}
	CatalogSource ModelCatalog
	Dest          interface{}
	CatalogDest   ModelCatalog
}

func (o *ModelAdapterOpts) Validate() error {
	var err error
	err = validateModel(o.Source)
	err = validateModel(o.Dest)
	sMtk, dMtk := modelT(o.Source).Kind(), modelT(o.Dest).Kind()
	if sMtk != dMtk {
		return fmt.Errorf("models must be of the same type. Received %s and %s",
			sMtk, dMtk)
	}
	return err
}

func (o *ModelAdapterOpts) Single() bool {
	return modelT(o.Dest).Kind() == reflect.Struct
}

// |||| BASE ADAPTER ||||

type ModelAdapter interface {
	Source() interface{}
	Dest() interface{}
	ExchangeToSource() error
	ExchangeToDest() error
}

func NewModelAdapter(opts *ModelAdapterOpts) ModelAdapter {
	if err := opts.Validate(); err != nil {
		log.Fatalln(err)
	}
	if opts.Single() {
		return NewSingleModelAdapter(opts)
	}
	return &MultiModelAdapter{opts: opts}
}

type ModelValues map[string]interface{}

// |||| MULTI MODEL ADAPTER ||||

type MultiModelAdapter struct {
	opts *ModelAdapterOpts
}

func (ma *MultiModelAdapter) exchange(toSource bool) error {
	var exchangeTo interface{}
	var exchangeFrom interface{}
	var catalog ModelCatalog
	if toSource {
		exchangeTo = ma.opts.Dest
		exchangeFrom = ma.opts.Source
		catalog = ma.opts.CatalogSource
	} else {
		exchangeTo = ma.opts.Source
		exchangeFrom = ma.opts.Dest
		catalog = ma.opts.CatalogDest
	}

	destRv := modelV(exchangeTo)
	sourceRv := modelV(exchangeFrom)
	for i := 0; i < destRv.Len(); i++ {
		destMv := destRv.Index(i).Interface()
		var sourceMv interface{}
		if i >= sourceRv.Len() {
			sourceMv = catalog.New(destMv)
		} else {
			sourceMv = sourceRv.Index(i).Interface()
		}
		opts := &ModelAdapterOpts{Source: sourceMv, Dest: destMv}
		sm := NewSingleModelAdapter(opts)
		if err := sm.ExchangeToSource(); err != nil {
			return err
		}
		sourceRv.Set(reflect.Append(sourceRv, reflect.ValueOf(sm.Source())))
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
	return reflect.ValueOf(ma.opts.Source).Elem().Interface()
}

func (ma *MultiModelAdapter) Dest() interface{} {
	return reflect.ValueOf(ma.opts.Dest).Elem().Interface()
}

// |||| MODEL ADAPTER ||||

type SingleModelAdapter struct {
	opts     *ModelAdapterOpts
	sourceAm *adaptedModel
	destAm   *adaptedModel
}

func NewSingleModelAdapter(opts *ModelAdapterOpts) *SingleModelAdapter {
	return &SingleModelAdapter{
		opts:     opts,
		sourceAm: &adaptedModel{model: opts.Source},
		destAm:   &adaptedModel{model: opts.Dest},
	}
}

func (ma *SingleModelAdapter) Source() interface{} {
	return ma.sourceAm.model
}

func (ma *SingleModelAdapter) Dest() interface{} {
	return ma.destAm.model
}

func (ma *SingleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.BindVals(ma.destAm.MapVals(), ma.opts.CatalogSource)
}

func (ma *SingleModelAdapter) ExchangeToDest() error {
	return ma.destAm.BindVals(ma.sourceAm.MapVals(), ma.opts.CatalogDest)
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	model interface{}
}

// BindVals binds a set of ModelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) BindVals(mv ModelValues, ca ModelCatalog) error {
	/* ORDER OF OPERATIONS
	1. Get the adapter model value
	2. For each value in the values to bind
		1. get the corresponding field name
		2. get the reflect.value from the model values
		3. If the type of the value adn the type of hte field are the same
		set the field to the value and call it a day
			1. If it isn't
				1. We need to check if the value is a struct or a slice
				2. if it isn't, throw an error
				3. if it is, attempt to adapt the model


	*/
	dv := modelV(mw.model)
	for k, rv := range mv {
		field := dv.FieldByName(k)
		val := reflect.ValueOf(rv)
		//if !val.IsValid() {
		//	continue
		//}
		if !field.CanSet() {
			continue
			//return fmt.Errorf("unsettable key %s in vals while binding to %val ",
			//	k, dv.Type())
		}
		vt, ft := val.Type(), field.Type()
		if vt != ft {
			valModelKind := modelT(val.Interface()).Kind()
			if valModelKind == reflect.Slice {
				valModelVal := modelV(val.Interface())
				if valModelVal.Len() == 0 {
					continue
				}
				o := &ModelAdapterOpts{
					Source:      val.Interface(),
					Dest:        field.Addr().Interface(),
					CatalogDest: ca,
				}
				ma := NewModelAdapter(o)
				if err := ma.ExchangeToDest(); err != nil {
					log.Fatalln(err)
				}
				val = reflect.ValueOf(ma.Dest())
			} else if valModelKind == reflect.Struct {
				valModelVal := modelV(val.Interface())
				if !valModelVal.IsValid() {
					continue
				}
				source := val.Interface()
				dest := reflect.New(field.Type().Elem()).Interface()
				o := &ModelAdapterOpts{
					Source: source,
					Dest:   dest,
				}
				ma := NewModelAdapter(o)
				if err := ma.ExchangeToDest(); err != nil {
					log.Fatalln(err)
				}
				val = reflect.ValueOf(ma.Dest())
			} else {
				return fmt.Errorf("(%s) invalid type %v for field '%s' with type %v "+
					"this is a no-op", dv.Type(), vt, k, ft)
			}
		}
		field.Set(val)
	}
	return nil
}

// MapVals maps adaptedModel fields to ModelValues.
func (mw *adaptedModel) MapVals() ModelValues {
	var mv = ModelValues{}
	dv := modelV(mw.model)
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		if f.Type().Kind() == reflect.Slice {
			mv[t.Name] = f.Addr().Interface()
		} else {
			mv[t.Name] = f.Interface()
		}
	}
	return mv
}

// |||| UTILITIES ||||

// || TYPE AND VALUE GETTING ||
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
	if ctk != reflect.Pointer {
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
