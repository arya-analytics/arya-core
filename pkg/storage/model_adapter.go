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
	//log.Info("POINTERS ", sourcePtr, destPtr)
	//log.Info("POINTER TYPES ", reflect.TypeOf(sourcePtr), reflect.TypeOf(destPtr))
	sourceRfl, destRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	if err := sourceRfl.Validate(); err != nil {
		return nil, err
	}
	if err := destRfl.Validate(); err != nil {
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
	var toAppend []*model.Reflect
	for i := 0; i < from.ChainValue().Len(); i++ {
		var toChainItem interface{}
		ap := true
		if i >= to.ChainValue().Len() {
			toChainItem = to.NewModel().Pointer()
		} else {
			ap = false
			toChainItem = to.ChainValueByIndex(i).Pointer()
		}
		maX, err := NewModelAdapter(from.ChainValueByIndex(i).Pointer(), toChainItem)
		if err != nil {
			return err
		}
		if err := maX.ExchangeToDest(); err != nil {
			return err
		}
		if ap {
			toAppend = append(toAppend, maX.Dest())
		}
	}
	for _, v := range toAppend {
		to.ChainAppend(v)
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
	for key, rawVal := range mv {
		field := mw.refl.Value().FieldByName(key)
		val := reflect.ValueOf(rawVal)
		if val.IsZero() || !field.IsValid() {
			continue
		}
		if val.Type() != field.Type() {
			var fieldRef *model.Reflect
			if field.Type().Kind() == reflect.Ptr {
				fieldRef = model.NewReflect(field.Interface())
			} else {
				fp := reflect.New(field.Type())
				fp.Elem().Set(field)
				fieldRef = model.NewReflect(fp.Interface())
			}

			if err := fieldRef.Validate(); err != nil {
				return NewError(ErrTypeInvalidField)
			}
			// The first thing we do is check if it's a nested refl we need to adapt
			var fieldPointer interface{}
			rawPointer := rawVal
			if fieldRef.IsChain() {
				rp := reflect.New(val.Type())
				rp.Elem().Set(val)
				rawPointer = rp.Interface()
				fieldPointer = fieldRef.Pointer()
			} else if fieldRef.IsStruct() {
				if field.IsNil() {
					fieldPointer = fieldRef.NewModel().Pointer()
				} else {
					fieldPointer = field
				}
			} else {
				return NewError(ErrTypeInvalidField)
			}

			rawRef := model.NewReflect(rawPointer)
			if err := rawRef.Validate(); err != nil {
				return NewError(ErrTypeInvalidField)
			}

			ma, err := NewModelAdapter(rawPointer, fieldPointer)
			if err != nil {
				return err
			}
			if err := ma.ExchangeToDest(); err != nil {
				return err
			}
			val = ma.Dest().ValueForSet()
		}
		field.Set(val)
	}
	return nil
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	var mv = modelValues{}
	for i := 0; i < mw.refl.Value().NumField(); i++ {
		f := mw.refl.Value().Field(i)
		t := mw.refl.Type().Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}
