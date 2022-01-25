package storage

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

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

func (mc ModelCatalog) InCatalog(modelPtr interface{}) bool {
	refM := model.NewReflect(modelPtr)
	for _, cm := range mc {
		refCm := model.NewReflect(cm)
		if err := refCm.Validate(); err != nil {
			panic(err)
		}
		if refM.Type().Name() == refCm.Type().Name() {
			return true
		}
	}
	return false
}

type modelValues map[string]interface{}

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

func (ma *ModelAdapter) ExchangeToDest() error {
	return ma.exchange(ma.Dest(), ma.Source())
}

func (ma *ModelAdapter) exchange(to *model.Reflect, from *model.Reflect) error {
	var pErr error
	from.ForEach(func(nRfl *model.Reflect, i int) {
		fromAm := &adaptedModel{rfl: nRfl}
		toRfl := to
		if i != -1 {
			if i >= to.ChainValue().Len() {
				toRfl = to.NewModel()
				to.ChainAppend(toRfl)
			} else {
				toRfl = to.ChainValueByIndex(i)
			}
		}
		toAm := &adaptedModel{rfl: toRfl}
		if err := toAm.bindVals(fromAm.mapVals()); err != nil {
			pErr = err
		}
	})
	return pErr
}

// |||| ADAPTED MODEL |||||

type adaptedModel struct {
	rfl *model.Reflect
}

// bindVals binds a set of modelValues to the adaptedModel fields.
// Returns an error for invalid / non-existent keys and invalid types.
func (mw *adaptedModel) bindVals(mv modelValues) error {
	for key, rv := range mv {
		fld := mw.rfl.Value().FieldByName(key)
		v := reflect.ValueOf(rv)
		if v.IsZero() || !fld.IsValid() {
			continue
		}
		if v.Type() != fld.Type() {
			if fld.Type().Kind() == reflect.Interface {
				impl := v.Type().Implements(fld.Type())
				if !impl {
					panic("doesn't implement interface")
				}
			} else {
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
	for i := 0; i < mw.rfl.Value().NumField(); i++ {
		f := mw.rfl.Value().Field(i)
		t := mw.rfl.Type().Field(i)
		mv[t.Name] = f.Interface()
	}
	return mv
}
