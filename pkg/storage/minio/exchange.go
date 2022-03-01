package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"reflect"
)

type dataValue struct {
	PK   model.PK
	Data *telem.ChunkData
}

type dataValueChain []*dataValue

type exchange struct {
	*model.Exchange
}

func newWrappedExchange(sma *model.Exchange) *exchange {
	return &exchange{sma}
}

func (m *exchange) bucket() string {
	return caseconv.PascalToKebab(m.Dest.Type().Name())
}

func (m *exchange) dataVals() dataValueChain {
	var c dataValueChain
	m.Dest.ForEach(func(rfl *model.Reflect, i int) {
		data := rfl.StructFieldByRole("telemChunkData")
		c = append(c, &dataValue{PK: rfl.PK(), Data: data.Interface().(*telem.ChunkData)})
	})
	return c
}

func (m *exchange) bindDataVals(dvc dataValueChain) {
	for _, dv := range dvc {
		rfl, ok := m.Dest.ValueByPK(dv.PK)
		if !ok {
			if m.Dest.IsChain() {
				newRfl := m.Dest.NewStruct()
				newRfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
				newRfl.StructFieldByName("ID").Set(dv.PK.Value())
				m.Dest.ChainAppend(newRfl)
			} else {
				if !m.Dest.PKField().IsZero() {
					panic("object store meta data mismatch")
				}
				m.Dest.StructFieldByName("ID").Set(dv.PK.Value())
				m.Dest.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
			}
		} else {
			rfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
		}
	}
}
