package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

const (
	idKey   = "ID"
	dataKey = "Data"
)

type DataValue struct {
	ID   string
	Data storage.Object
}

type DataValueChain []*DataValue

func (dvc DataValueChain) Retrieve(ID string) *DataValue {
	for _, dv := range dvc {
		if dv.ID == ID {
			return dv
		}
	}
	return nil
}

func (dvc DataValueChain) Contains(ID string) (e bool) {
	for _, dv := range dvc {
		if dv.ID == ID {
			e = true
		}
	}
	return e
}

type ModelWrapper struct {
	rfl *model.Reflect
}

func (m *ModelWrapper) Bucket() string {
	return m.rfl.Type().Name()
}

func (m *ModelWrapper) DataVals() DataValueChain {
	var c DataValueChain
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		val := rfl.Value()
		id := rfl.IDField().String()
		data := val.FieldByName(dataKey)
		c = append(c, &DataValue{
			ID:   id,
			Data: data.Interface().(storage.Object)})
	})
	return c
}

func (m *ModelWrapper) BindDataVals(dvc *DataValueChain) {
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		val := rfl.Value()
		id := rfl.IDField().String()
		data := val.FieldByName(dataKey)
		if dvc.Contains(id) {
			data.Set(reflect.ValueOf(dvc.Retrieve(id)))
		}
	})
}
