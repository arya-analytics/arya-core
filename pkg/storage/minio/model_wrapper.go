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

type ModelWrapper struct {
	rfl *model.Reflect
}

func (m *ModelWrapper) Bucket() string {
	return caseconv.PascalToKebab(m.rfl.Type().Name())
}

func (m *ModelWrapper) DataVals() DataValueChain {
	var c DataValueChain
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		val := rfl.StructValue()
		data := val.FieldByName(dataKey)
		c = append(c, &DataValue{
			PK:   rfl.PK(),
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
				newRfl := m.rfl.NewStruct()
				newRfl.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
				newRfl.StructValue().FieldByName("ID").Set(dv.PK.Value())
				m.rfl.ChainAppend(newRfl)
			} else {
				if !m.rfl.PKField().IsZero() {
					panic("object store meta data mismatch")
				}
				m.rfl.StructValue().FieldByName("ID").Set(dv.PK.Value())
				m.rfl.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
			}
		} else {
			rfl.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
		}
	}
}
