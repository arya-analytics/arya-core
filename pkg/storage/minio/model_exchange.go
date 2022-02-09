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

type dataValue struct {
	PK   model.PK
	Data storage.Object
}

type dataValueChain []*dataValue

type modelExchange struct {
	*model.Exchange
}

func newWrappedModelExchange(sma *model.Exchange) *modelExchange {
	return &modelExchange{sma}
}

func (m *modelExchange) Bucket() string {
	return caseconv.PascalToKebab(m.Dest.Type().Name())
}

func (m *modelExchange) DataVals() dataValueChain {
	var c dataValueChain
	m.Dest.ForEach(func(rfl *model.Reflect, i int) {
		val := rfl.StructValue()
		data := val.FieldByName(dataKey)
		c = append(c, &dataValue{
			PK:   rfl.PK(),
			Data: data.Interface().(storage.Object),
		})
	})
	return c
}

func (m *modelExchange) BindDataVals(dvc dataValueChain) {
	for _, dv := range dvc {
		rfl, ok := m.Dest.ValueByPK(dv.PK)
		if !ok {
			if m.Dest.IsChain() {
				newRfl := m.Dest.NewStruct()
				newRfl.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
				newRfl.StructValue().FieldByName("ID").Set(dv.PK.Value())
				m.Dest.ChainAppend(newRfl)
			} else {
				if !m.Dest.PKField().IsZero() {
					panic("object store meta data mismatch")
				}
				m.Dest.StructValue().FieldByName("ID").Set(dv.PK.Value())
				m.Dest.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
			}
		} else {
			rfl.StructValue().FieldByName(dataKey).Set(reflect.ValueOf(dv.Data))
		}
	}
}
