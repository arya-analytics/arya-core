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

/// |||| CATALOG ||||

type ModelCatalog []reflect.Type

func (mc ModelCatalog) New(m interface{}) interface{} {
	mName := modelT(m).Name()
	for _, cm := range mc {
		if cm.Name() == mName {
			return reflect.New(cm).Interface()
		}
	}
	log.Fatalf("Model %s could not be found in catalog. This is an no-op.", mName)
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
	// Validation checks
	// 1. Model and Source are pointers to a struct or a slice that can be set
	err = validateModel(o.Source)
	err = validateModel(o.Dest)
	sMtk, dMtk := modelT(o.Source).Kind(), modelT(o.Dest).Kind()
	if sMtk != dMtk {
		return fmt.Errorf("models must be of the same type. Received %s and %s",
			sMtk, dMtk)
	}
	return err
}

func (o *ModelAdapterOpts) single() bool {
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
	if opts.single() {
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
		catalog = ma.opts.CatalogDest
	} else {
		exchangeTo = ma.opts.Source
		exchangeFrom = ma.opts.Dest
		catalog = ma.opts.CatalogSource
	}

	destRv := modelV(exchangeTo)
	sourceRv := modelV(exchangeFrom)
	for i := 0; i < destRv.Len(); i++ {
		destMv := destRv.Index(i).Interface()
		var sourceMv interface{}
		if i >= sourceRv.Len() {
			log.Warn(destMv)
			sourceMv = catalog.New(destMv)
		} else {
			sourceMv = sourceRv.Index(i).Interface()
		}
		opts := &ModelAdapterOpts{Source: sourceMv, Dest: destMv}
		ma := NewSingleModelAdapter(opts)
		if err := ma.ExchangeToSource(); err != nil {
			return err
		}
		sourceRv.Set(reflect.Append(sourceRv, reflect.ValueOf(ma.Source())))
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
	sourceAm *AdaptedModel
	destAm   *AdaptedModel
}

func NewSingleModelAdapter(opts *ModelAdapterOpts) *SingleModelAdapter {
	return &SingleModelAdapter{
		opts:     opts,
		sourceAm: &AdaptedModel{Model: opts.Source},
		destAm:   &AdaptedModel{Model: opts.Dest},
	}
}

func (ma *SingleModelAdapter) Source() interface{} {
	return ma.sourceAm.Model
}

func (ma *SingleModelAdapter) Dest() interface{} {
	return ma.destAm.Model
}

func (ma *SingleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.BindVals(ma.destAm.MapVals())
}

func (ma *SingleModelAdapter) ExchangeToDest() error {
	return ma.destAm.BindVals(ma.sourceAm.MapVals())
}

// |||| ADAPTED MODEL |||||

type AdaptedModel struct {
	Model interface{}
}

// BindVals binds a set of ModelValues to the AdaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *AdaptedModel) BindVals(mv ModelValues) error {
	dv := modelV(mw.Model)
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
	dv := modelV(mw.Model)
	for i := 0; i < dv.NumField(); i++ {
		t := dv.Type().Field(i)
		f := dv.Field(i)
		mv[t.Name] = f.Interface()
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
