package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"reflect"
)

const (
	roleStamp = "tsStamp"
	roleKey   = "tsKey"
	roleValue = "tsValue"
)

type exchange struct {
	*model.Exchange
}

func newWrappedExchange(sma *model.Exchange) *exchange {
	return &exchange{sma}
}

func (m *exchange) samples() (samples []timeseries.Sample) {
	m.Dest().ForEach(func(rfl *model.Reflect, i int) {
		samples = append(samples, m.newSampleFromRFL(rfl))
	})
	return samples
}

func (m *exchange) seriesNames() (names []string) {
	m.Dest().ForEach(func(rfl *model.Reflect, i int) {
		pk := model.NewPK(keyField(rfl).Interface())
		names = append(names, pk.String())
	})
	return names
}

func (m *exchange) bindRes(key string, res interface{}) error {
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

func (m *exchange) newSampleFromRFL(rfl *model.Reflect) timeseries.Sample {
	return timeseries.Sample{
		Key:       model.NewPK(keyField(rfl).Interface()).String(),
		Timestamp: stampField(rfl).Interface().(telem.TimeStamp),
		Value:     valueField(rfl).Interface().(float64),
	}
}

func (m *exchange) appendSample(sample timeseries.Sample) {
	newRfl := m.Dest().NewStruct()
	m.setFields(newRfl, sample)
	m.Dest().ChainAppend(newRfl)
}

func (m *exchange) setFields(rfl *model.Reflect, sample timeseries.Sample) {
	kf := keyField(rfl)
	kf.Set(convertKeyString(kf.Type(), sample.Key))
	stampField(rfl).Set(reflect.ValueOf(sample.Timestamp))
	valueField(rfl).Set(reflect.ValueOf(sample.Value))
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

// |||| UTILITIES |||

func fieldByRole(rfl *model.Reflect, role string) reflect.Value {
	return rfl.StructFieldByRole(role)
}

func keyField(rfl *model.Reflect) reflect.Value {
	return fieldByRole(rfl, roleKey)
}

func valueField(rfl *model.Reflect) reflect.Value {
	return fieldByRole(rfl, roleValue)
}

func stampField(rfl *model.Reflect) reflect.Value {
	return fieldByRole(rfl, roleStamp)
}
