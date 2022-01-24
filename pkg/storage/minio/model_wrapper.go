package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

const (
	dataKey = "Data"
)

type DataValue struct {
	PK   model.PK
	Data storage.Object
}

type DataValueChain []*DataValue

func (dvc DataValueChain) Retrieve(pk model.PK) *DataValue {
	for _, dv := range dvc {
		if dv.PK.Equals(pk) {
			return dv
		}
	}
	return nil
}

func (dvc DataValueChain) Contains(pk model.PK) bool {
	for _, dv := range dvc {
		if dv.PK.Equals(pk) {
			return true
		}
	}
	return false
}

type ModelWrapper struct {
	rfl *model.Reflect
}

func (m *ModelWrapper) Bucket() string {
	return caseconv.PascalToKebab(m.rfl.Type().Name())
}

func (m *ModelWrapper) DataVals() DataValueChain {
	var c DataValueChain
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		val := rfl.Value()
		data := val.FieldByName(dataKey)
		c = append(c, &DataValue{
			PK:   rfl.PKField(),
			Data: data.Interface().(storage.Object),
		})
	})
	return c
}

func (m *ModelWrapper) BindDataVals(dvc DataValueChain) {
	for _, dv := range dvc {
		rfl, ok := m.rfl.ValueByPK(dv.PK)
		if !ok {
			if m.rfl.IsChain() {
				newRfl := m.rfl.NewModel()
				newRfl.Value().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
				m.rfl.ChainAppend(newRfl)
			} else {
				if !m.rfl.PKField().IsZero() {
					panic("object store meta data mismatch")
				}
				m.rfl.Value().FieldByName(model.KeyPK).Set(dv.PK.Value())
				m.rfl.Value().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
			}
		} else {
			rfl.Value().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
		}
	}
}
