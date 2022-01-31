package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
)

const (
	roleStamp = "tsStamp"
	roleKey   = "tsKey"
	roleValue = "tsValue"
)

type modelAdapter struct {
	*storage.ModelAdapter
}

func newWrappedModelAdapter(sma *storage.ModelAdapter) *modelAdapter {
	return &modelAdapter{sma}
}

func (m *modelAdapter) samples() (samples []timeseries.Sample) {
	m.Dest().ForEach(func(rfl *model.Reflect, i int) {
		samples = append(samples, m.newSampleFromRFL(rfl))
	})
	return samples
}

func (m *modelAdapter) seriesNames() (names []string) {
	m.Dest().ForEach(func(rfl *model.Reflect, i int) {
		fld := rfl.StructFieldByRole(roleKey)
		pk := model.NewPK(fld.Interface())
		names = append(names, pk.String())
	})
	return names
}

func (m *modelAdapter) bindRes(key string, res interface{}) error {
	resVal := reflect.ValueOf(res)
	if resVal.Type().Kind() != reflect.Slice {
		panic("received unknown response from cache")
	}
	if resVal.Len() == 0 {
		return nil
	}
	if resVal.Index(0).Elem().Type().Kind() == reflect.Slice {
		for i := 0; i < resVal.Len(); i++ {
			resItemVal := resVal.Index(i)
			sample, err := timeseries.NewSampleFromRes(key, resItemVal.Interface())
			if err != nil {
				return err
			}
			if m.Dest().IsChain() {
				m.appendSample(sample)
			} else {
				panic("can't bind multiple result values to a non-chain")
			}
		}
	} else {
		sample, err := timeseries.NewSampleFromRes(key, res)
		if err != nil {
			return err
		}
		if m.Dest().IsChain() {
			m.appendSample(sample)
		} else {
			m.setFields(m.Dest(), sample)
		}
	}
	return nil
}

func (m *modelAdapter) newSampleFromRFL(rfl *model.Reflect) timeseries.Sample {
	return timeseries.Sample{
		Key:       model.NewPK(rfl.StructFieldByRole(roleKey).Interface()).String(),
		Timestamp: rfl.StructFieldByRole(roleStamp).Interface().(int64),
		Value:     rfl.StructFieldByRole(roleValue).Interface().(float64),
	}
}

func (m *modelAdapter) appendSample(sample timeseries.Sample) {
	newRfl := m.Dest().NewStruct()
	m.setFields(newRfl, sample)
	m.Dest().ChainAppend(newRfl)
}

func (m *modelAdapter) setFields(rfl *model.Reflect, sample timeseries.Sample) {
	kf := rfl.StructFieldByRole(roleKey)
	kf.Set(convertKeyString(kf.Type(), sample.Key))
	rfl.StructFieldByRole(roleStamp).Set(reflect.ValueOf(sample.Timestamp))
	rfl.StructFieldByRole(roleValue).Set(reflect.ValueOf(sample.Value))
}

func convertKeyString(t reflect.Type, keyString string) reflect.Value {
	switch t {
	case reflect.TypeOf(uuid.UUID{}):
		id, err := uuid.Parse(keyString)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(id)
	}
	panic("received unexpected type")
}
