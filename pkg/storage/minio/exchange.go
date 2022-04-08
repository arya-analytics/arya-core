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

type reflectMinio struct {
	*model.Reflect
}

func wrapReflect(rfl *model.Reflect) *reflectMinio {
	return &reflectMinio{rfl}
}

func (m *reflectMinio) bucket() string {
	return caseconv.PascalToKebab(m.Type().Name())
}

func (m *reflectMinio) dataValues() (dvc []data) {
	m.ForEach(func(rfl *model.Reflect, i int) {
		d := rfl.StructFieldByRole("telemChunkData").Interface().(*telem.ChunkData)
		dvc = append(dvc, data{
			PK:   rfl.PK(),
			Data: d,
		})
	})
	return dvc
}

func (m *reflectMinio) bindDataVals(dvc []data) {
	for _, dv := range dvc {
		rfl, ok := m.ValueByPK(dv.PK)
		if ok {
			rfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
		} else {
			nRfl := m.Reflect
			if m.IsChain() {
				nRfl = m.NewStruct()
			}
			nRfl.StructFieldByRole("telemChunkData").Set(reflect.ValueOf(dv.Data))
			nRfl.StructFieldByName("ID").Set(dv.PK.Value())
			if m.IsChain() {
				m.ChainAppend(nRfl)
			}
		}
	}
}
