package redis

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
)

const (
	structTagCacheCat = "cache"
	// || ROLES ||
	structTagRoleKey   = "role"
	structTagStampRole = "tsStamp"
	structTagKeyRole   = "tsKey"
	structTagValRole   = "tsValue"
)

type tsModelWrapper struct {
	rfl *model.Reflect
}

func (m *tsModelWrapper) samples() (samples []timeseries.Sample) {
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		samples = append(samples, m.newSampleFromRFL(rfl))
	})
	return samples
}

func (m *tsModelWrapper) seriesNames() (names []string) {
	m.rfl.ForEach(func(rfl *model.Reflect, i int) {
		fld := m.retrieveFieldByRole(rfl, structTagKeyRole).Interface()
		pk := model.NewPK(fld)
		names = append(names, pk.String())
	})
	return names
}

func (m *tsModelWrapper) bindRes(key string, res interface{}) error {
	resVal := reflect.ValueOf(res)
	if resVal.Type().Kind() != reflect.Slice {
		panic("received unknown response from cache")
	}
	if resVal.Len() == 0 {
		return nil
	}
	if resVal.Index(0).Type().Kind() == reflect.Slice {
		for i := 0; i < resVal.Len(); i++ {
			resItemVal := resVal.Index(i)
			sample, err := timeseries.NewSampleFromRes(key, resItemVal.Interface())
			if err != nil {
				return err
			}
			if m.rfl.IsChain() {
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
		if m.rfl.IsChain() {
			m.appendSample(sample)
		} else {
			m.setFields(m.rfl, sample)
		}
	}
	return nil
}

func (m *tsModelWrapper) newSampleFromRFL(rfl *model.Reflect) timeseries.Sample {
	sfn, kfn, vfn := m.fieldNames()
	stamp := rfl.Value().FieldByName(sfn).Interface()
	key := rfl.Value().FieldByName(kfn).Interface()
	pk := model.NewPK(key)
	value := rfl.Value().FieldByName(vfn).Interface()
	return timeseries.Sample{
		Key:       pk.String(),
		Timestamp: stamp.(int64),
		Value:     value.(float64),
	}
}

func (m *tsModelWrapper) appendSample(sample timeseries.Sample) {
	newRfl := m.rfl.NewModel()
	m.setFields(newRfl, sample)
	m.rfl.ChainAppend(newRfl)
}

func (m *tsModelWrapper) fieldNames() (sfn string, kfn string, vfn string) {
	sfn = m.retrieveFieldNameByRole(structTagStampRole)
	kfn = m.retrieveFieldNameByRole(structTagKeyRole)
	vfn = m.retrieveFieldNameByRole(structTagValRole)
	return sfn, kfn, vfn
}

func (m *tsModelWrapper) retrieveFieldByRole(rfl *model.Reflect,
	role string) reflect.Value {
	fldName := m.retrieveFieldNameByRole(role)
	return rfl.Value().FieldByName(fldName)
}

func (m *tsModelWrapper) setFields(rfl *model.Reflect, sample timeseries.Sample) {
	sfn, kfn, vfn := m.fieldNames()
	kf := m.retrieveFieldByRole(rfl, structTagKeyRole)
	rfl.Value().FieldByName(kfn).Set(convertKeyString(kf.Type(), sample.Key))
	rfl.Value().FieldByName(sfn).Set(reflect.ValueOf(sample.Timestamp))
	rfl.Value().FieldByName(vfn).Set(reflect.ValueOf(sample.Value))
}

func (m *tsModelWrapper) retrieveFieldNameByRole(role string) string {
	t, ok := m.rfl.Tags().Retrieve(structTagCacheCat, structTagRoleKey, role)
	if !ok {
		panic(fmt.Sprintf("couldn't retrieve role %s from model", role))
	}
	return t.FldName
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
