package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"reflect"
)

const (
	dataKey = "Telem"
)

type dataValue struct {
	PK   model.PK
	Data *telem.Bulk
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
			Data: data.Interface().(*telem.Bulk),
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
				newRfl.StructFieldByRole("bulkTelem").Set(reflect.ValueOf(dv.Data))
				newRfl.StructFieldByName("ID").Set(dv.PK.Value())
				m.Dest.ChainAppend(newRfl)
			} else {
				if !m.Dest.PKField().IsZero() {
					panic("object store meta data mismatch")
				}
				m.Dest.StructFieldByName("ID").Set(dv.PK.Value())
				m.Dest.StructFieldByRole("bulkTelem").Set(reflect.ValueOf(dv.Data))
			}
		} else {
			rfl.StructFieldByRole("bulkTelem").Set(reflect.ValueOf(dv.Data))
		}
	}
}
