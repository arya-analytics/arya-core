package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"reflect"
)

type data struct {
	PK   model.PK
	Data *telem.ChunkData
}

type exchange struct {
	*model.Exchange
}

func wrapExchange(sma *model.Exchange) *exchange {
	return &exchange{sma}
}

func (m *exchange) bucket() string {
	return caseconv.PascalToKebab(m.Dest().Type().Name())
}

func (m *exchange) dataVals() (dvc []data) {
	m.Dest().ForEach(func(rfl *model.Reflect, i int) {
		dvc = append(dvc, data{
			PK:   rfl.PK(),
			Data: rfl.StructFieldByRole("telemChunkData").Interface().(*telem.ChunkData),
		})
	})
	return dvc
}

func (m *exchange) bindDataVals(dvc []data) {
	for _, dv := range dvc {
		rfl, ok := m.Dest().ValueByPK(dv.PK)
		if ok {
			rfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
		} else {
			nRfl := m.Dest()
			if m.Dest().IsChain() {
				nRfl = m.Dest().NewStruct()
			}
			nRfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
			nRfl.StructFieldByName("ID").Set(dv.PK.Value())
			if m.Dest().IsChain() {
				m.Dest().ChainAppend(nRfl)
			}
		}
	}
}
