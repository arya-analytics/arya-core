package storage

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

/// |||| CATALOG ||||

type ModelCatalog []interface{}

func (mc ModelCatalog) New(modelPtr interface{}) interface{} {
	refM := model.NewReflect(modelPtr)
	for _, cm := range mc {
		refCm := model.NewReflect(cm)
		if err := refCm.Validate(); err != nil {
			panic(err)
		}
		if refM.Type().Name() == refCm.Type().Name() {
			if refM.IsChain() {
				return refCm.NewChain().Pointer()
			}
			return refCm.NewModel().Pointer()
		}
	}
	panic(fmt.Sprintf("model %s could not be found in catalog", refM.Type().Name()))
}

// |||| BASE ADAPTER ||||

type modelValues map[string]interface{}

// |||| MODEL ADAPTER ||||

type ModelAdapter struct {
	sourceRfl *model.Reflect
	destRfl   *model.Reflect
}

func NewModelAdapter(sourcePtr interface{}, destPtr interface{}) *ModelAdapter {
	sourceRfl, destRfl := model.NewReflect(sourcePtr), model.NewReflect(destPtr)
	if err := sourceRfl.Validate(); err != nil {
		panic(err)
	}
	if err := destRfl.Validate(); err != nil {
		panic(err)
	}
	if sourceRfl.RawType().Kind() != destRfl.RawType().Kind() {
		panic("model adapter received model and chain. source and dest must be equal")
	}
	return &ModelAdapter{sourceRfl, destRfl}
}

func (ma *ModelAdapter) Source() *model.Reflect {
	return ma.sourceRfl
}

func (ma *ModelAdapter) Dest() *model.Reflect {
	return ma.destRfl
}

func (ma *ModelAdapter) ExchangeToSource() error {
	return ma.exchange(ma.Source(), ma.Dest())
}

func (ma *ModelAdapter) exchange(to *model.Reflect, from *model.Reflect) error {
	var pErr error
	from.ForEach(func(nRfl *model.Reflect, i int) {
		fromAm := &adaptedModel{refl: nRfl}
		var toAm *adaptedModel
		if i == -1 {
			toAm = &adaptedModel{refl: to}
		} else {
			if i >= to.ChainValue().Len() {
				newM := to.NewModel()
				toAm = &adaptedModel{refl: newM}
				to.ChainAppend(newM)
			} else {
				toAm = &adaptedModel{
					refl: to.ChainValueByIndex(i),
				}
			}
		}
		err := toAm.bindVals(fromAm.mapVals())
		if err != nil {
			pErr = err
		}
	})
	return pErr
}

func (ma *ModelAdapter) ExchangeToDest() error {
	return ma.exchange(ma.Dest(), ma.Source())
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	refl *model.Reflect
}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) error {
	for key, rv := range mv {
		fld := mw.refl.Value().FieldByName(key)
		v := reflect.ValueOf(rv)
		if v.IsZero() || !fld.IsValid() {
			continue
		}
		if v.Type() != fld.Type() {
			fldRfl, err := mw.newValidatedRfl(fld.Interface())
			if err != nil {
				return err
			}
			vRfl, err := mw.newValidatedRfl(rv)
			if err != nil {
				return err
			}
			vPtr := vRfl.Pointer()
			fldPtr := fldRfl.Pointer()
			if fldRfl.IsStruct() && fld.IsNil() {
				fldPtr = fldRfl.NewModel().Pointer()
			}
			ma := NewModelAdapter(vPtr, fldPtr)
			if err := ma.ExchangeToDest(); err != nil {
				return err
			}
			v = ma.Dest().ValueForSet()
		}
		fld.Set(v)
	}
	return nil
}

func (mw *adaptedModel) newValidatedRfl(v interface{}) (*model.Reflect, error) {
	rfl := model.NewReflect(v)
	if !rfl.IsPointer() {
		rfl = rfl.NewPointer()
	}
	if err := rfl.Validate(); err != nil {
		return nil, NewError(ErrTypeInvalidField)
	}
	return rfl, nil
}

// mapVals maps adaptedModel fields to modelValues.
func (mw *adaptedModel) mapVals() modelValues {
	mv := modelValues{}
	for i := 0; i < mw.refl.Value().NumField(); i++ {
		f := mw.refl.Value().Field(i)
		t := mw.refl.Type().Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}
