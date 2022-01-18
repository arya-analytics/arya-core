package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

// TODO: Move general functions to utilities
// TODO: Figure out how to simplify type getting system
// TODO: Document APIs

/// |||| CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(modelPtr interface{}) interface{} {
	refM := model.NewReflect(modelPtr)
	for _, cm := range mc {
		refCm := model.NewReflect(cm)
		if refM.Name() == refCm.Name() {
			if refM.IsChain() {
				return refM.NewChain().Pointer()
			}
			return refM.NewModel().Pointer()
		}
	}
	log.Fatalf("modelReflect %s could not be found in catalog. This is an no-op.", refM.Name())
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

func NewModelAdapter(sourcePtr interface{}, destPtr interface{}) (ModelAdapter, error) {
	refSource := model.NewReflect(sourcePtr)
	if err := refSource.Validate(); err != nil {
		return nil, err
	}
	refDest := model.NewReflect(destPtr)
	if err := refDest.Validate(); err != nil {
		return nil, err
	}
	if refSource.Type().Kind() != refDest.Type().Kind() {
		return nil, model.NewError(model.ErrTypeIncompatibleModels)
	}
	if refSource.IsChain() {
		return newSingleModelAdapter(refSource, refDest), nil
	}
	return &chainModelAdapter{refSource, refDest}, nil
}

type modelValues map[string]interface{}

// |||| MULTI MODEL ADAPTER ||||

type chainModelAdapter struct {
	source *model.Reflect
	dest   *model.Reflect
}

func (ma *chainModelAdapter) exchange(toSource bool) error {
	var to *model.Reflect
	var from *model.Reflect
	if toSource {
		to, from = ma.source, ma.dest
	} else {
		to, from = ma.dest, ma.source
	}
	for i := 0; i < from.ChainValue().Len(); i++ {
		var toChainItem interface{}
		if i >= to.ChainValue().Len() {
			tc := to.NewModel()
			toChainItem = tc.Pointer()
		} else {
			toChainItem = to.ValueIndex(i)
		}
		maX, err := NewModelAdapter(from.ValueIndex(i).Interface(), toChainItem)
		if err != nil {
			return err
		}
		if err := maX.ExchangeToDest(); err != nil {
			return err
		}
		to.ChainAppend(maX.DestValue())
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
	return ma.source.Pointer()
}

func (ma *chainModelAdapter) DestPointer() interface{} {
	return ma.source.Pointer()
}

func (ma *chainModelAdapter) DestValue() reflect.Value {
	return ma.dest.Value()
}

func (ma *chainModelAdapter) SourceValue() reflect.Value {
	return ma.source.Value()
}

// |||| MODEL ADAPTER ||||

type singleModelAdapter struct {
	sourceAm *adaptedModel
	destAm   *adaptedModel
}

func newSingleModelAdapter(source *model.Reflect, dest *model.Reflect) *singleModelAdapter {
	return &singleModelAdapter{
		sourceAm: &adaptedModel{modelReflect: source},
		destAm:   &adaptedModel{modelReflect: dest},
	}
}

func (ma *singleModelAdapter) SourcePointer() interface{} {
	return ma.sourceAm.modelReflect.Pointer()
}

func (ma *singleModelAdapter) DestPointer() interface{} {
	return ma.destAm.modelReflect.Pointer()
}

func (ma *singleModelAdapter) DestValue() reflect.Value {
	return ma.destAm.modelReflect.Value()
}

func (ma *singleModelAdapter) SourceValue() reflect.Value {
	return ma.sourceAm.modelReflect.Value()
}

func (ma *singleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.bindVals(ma.destAm.mapVals())
}

func (ma *singleModelAdapter) ExchangeToDest() error {
	return ma.destAm.bindVals(ma.sourceAm.mapVals())
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	modelReflect *model.Reflect
}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) error {
	for k, rv := range mv {
		field := mw.modelReflect.StructFieldByName(k)
		val := reflect.ValueOf(rv)
		if !field.CanSet() {
			continue
		}
		vt, ft := val.Type(), field.Type()
		if vt != ft {
			// The first thing we do is check if it's a nested modelReflect we need to adapt
			fieldRef := model.NewReflect(field.Interface())
			if err := fieldRef.Validate(); err != nil {
				return NewError(ErrTypeInvalidField)
			}
			ma, err := NewModelAdapter(val, field)
			if err != nil {
				return err
			}
			if err := ma.ExchangeToDest(); err != nil {
				return err
			}
			val = ma.DestValue()
		}
		field.Set(val)
	}
	return nil
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	var mv = modelValues{}
	for i := 0; i < mw.modelReflect.StructNumFields(); i++ {
		f := mw.modelReflect.StructFieldByIndex(i)
		t := mw.modelReflect.Type().Field(i)
		// Need to convert slices to addy because that's what NewModelAdapter expects.
		if f.Type().Kind() == reflect.Slice {
			mv[t.Name] = f.Addr().Interface()
		} else {
			mv[t.Name] = f.Interface()
		}
	}
	return mv
}
