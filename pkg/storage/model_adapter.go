package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

// TODO: Document APIs

/// |||| CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(modelPtr interface{}) interface{} {
	refM := model.NewReflect(modelPtr)
	for _, cm := range mc {
		refCm := model.NewReflect(cm)
		if err := refCm.Validate(); err != nil {
			log.Fatalln(err)
		}
		if refM.Type().Name() == refCm.Type().Name() {
			if refM.IsChain() {
				return refM.NewChain().Pointer()
			}
			return refM.NewModel().Pointer()
		}
	}
	log.Fatalf("model %s could not be found in catalog", refM.Type().Name())
	return nil
}

// |||| BASE ADAPTER ||||

type ModelAdapter interface {
	Source() *model.Reflect
	Dest() *model.Reflect
	ExchangeToSource() error
	ExchangeToDest() error
}

func NewModelAdapter(sourcePtr interface{}, destPtr interface{}) (ModelAdapter, error) {
	log.Info("POINTERS ", sourcePtr, destPtr)
	log.Info("POINTER TYPES ", reflect.TypeOf(sourcePtr), reflect.TypeOf(destPtr))
	sourceRfl, destRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	err := sourceRfl.Validate()
	err = destRfl.Validate()
	if err != nil {
		return nil, err
	}
	if sourceRfl.RawType().Kind() != destRfl.RawType().Kind() {
		return nil, model.NewError(model.ErrTypeIncompatibleModels)
	}
	if sourceRfl.IsStruct() {
		return newSingleModelAdapter(sourceRfl, destRfl), nil
	}
	return &chainModelAdapter{sourceRfl, destRfl}, nil
}

type modelValues map[string]interface{}

// |||| MULTI MODEL ADAPTER ||||

type chainModelAdapter struct {
	sourceRfl *model.Reflect
	destRfl   *model.Reflect
}

func (ma *chainModelAdapter) exchange(toSource bool) error {
	var to *model.Reflect
	var from *model.Reflect
	if toSource {
		to, from = ma.sourceRfl, ma.destRfl
	} else {
		to, from = ma.destRfl, ma.sourceRfl
	}
	for i := 0; i < from.ChainValue().Len(); i++ {
		var toChainItem interface{}
		if i >= to.ChainValue().Len() {
			toChainItem = to.NewModel().Pointer()
		} else {
			toChainItem = to.ChainValueByIndex(i)
		}
		maX, err := NewModelAdapter(from.ChainValueByIndex(i).Pointer(), toChainItem)
		if err != nil {
			return err
		}
		if err := maX.ExchangeToDest(); err != nil {
			return err
		}
		to.ChainAppend(maX.Dest())
	}
	return nil
}

func (ma *chainModelAdapter) ExchangeToSource() error {
	return ma.exchange(true)
}

func (ma *chainModelAdapter) ExchangeToDest() error {
	return ma.exchange(false)
}

func (ma *chainModelAdapter) Source() *model.Reflect {
	return ma.sourceRfl
}

func (ma *chainModelAdapter) Dest() *model.Reflect {
	return ma.destRfl
}

// |||| MODEL ADAPTER ||||

type singleModelAdapter struct {
	sourceAm *adaptedModel
	destAm   *adaptedModel
}

func newSingleModelAdapter(source *model.Reflect, dest *model.Reflect) *singleModelAdapter {
	return &singleModelAdapter{
		sourceAm: &adaptedModel{refl: source},
		destAm:   &adaptedModel{refl: dest},
	}
}

func (ma *singleModelAdapter) Source() *model.Reflect {
	return ma.sourceAm.refl
}

func (ma *singleModelAdapter) Dest() *model.Reflect {
	return ma.destAm.refl
}

func (ma *singleModelAdapter) ExchangeToSource() error {
	return ma.sourceAm.bindVals(ma.destAm.mapVals())
}

func (ma *singleModelAdapter) ExchangeToDest() error {
	return ma.destAm.bindVals(ma.sourceAm.mapVals())
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	refl *model.Reflect
}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) error {
	for k, rv := range mv {
		field := mw.refl.Value().FieldByName(k)
		val := reflect.ValueOf(rv)
		valType := reflect.TypeOf(rv)
		if field.Type().Kind() == reflect.Slice {
			valType = valType.Elem()
		}
		if !field.CanSet() {
			continue
		}
		log.Info(valType, field.Type())
		if valType != field.Type() {
			// The first thing we do is check if it's a nested refl we need to adapt
			fieldPointer := field.Interface()
			if field.Type().Kind() == reflect.Slice {
				fieldPointer = field.Addr().Interface()
			}
			fieldRef := model.NewReflect(fieldPointer)
			if err := fieldRef.Validate(); err != nil {
				return NewError(ErrTypeInvalidField)
			}
			if fieldRef.IsStruct() && val.IsNil() {
				continue
			}
			ma, err := NewModelAdapter(rv, fieldRef.NewRaw().Pointer())
			if err != nil {
				return err
			}
			if err := ma.ExchangeToDest(); err != nil {
				return err
			}
			val = ma.Dest().RawValue()
		}
		if field.Type().Kind() == reflect.Slice {
			val = val.Elem()
		}
		log.Warn(val)
		field.Set(val)
	}
	return nil
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	var mv = modelValues{}
	for i := 0; i < mw.refl.Value().NumField(); i++ {
		f := mw.refl.Value().Index(i)
		t := mw.refl.Type().Field(i)
		if f.Type().Kind() == reflect.Slice {
			mv[t.Name] = f.Addr().Interface()
		} else {
			mv[t.Name] = f.Interface()
		}
	}
	return mv
}
