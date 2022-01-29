package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type tsBaseQuery struct {
	pks      []interface{}
	modelRfl *model.Reflect
	baseQuery
}

func (tsb *tsBaseQuery) tsBaseModel(m interface{}) {
	tsb.modelRfl = model.NewReflect(m)
}

func (tsb *tsBaseQuery) tsBaseWherePk(pk interface{}) {
	tsb.pks = append(tsb.pks, pk)
}

func (tsb *tsBaseQuery) tsBaseWherePks(pks interface{}) {
	rv := reflect.ValueOf(pks)
	for i := 0; i < rv.Len(); i++ {
		tsb.tsBaseWherePk(rv.Index(i).Interface())
	}
}

func (tsb *tsBaseQuery) tsBaseCreateIndexes(ctx context.Context,
	pks []interface{}) error {
	log.Warn("Series not found in cache. Attempting to fix.")
	tag, ok := tsb.modelRfl.Tags().Retrieve("storage", "role", "index")
	if !ok {
		panic("couldn't get tag from model")
	}
	fld, ok := tsb.modelRfl.Type().FieldByName(tag.FldName)
	if !ok {
		panic("couldn't get field")
	}
	sm := catalog().NewFromType(fld.Type.Elem(), true)
	smRfl := model.NewReflect(sm)
	if sErr := tsb.storage.NewRetrieve().Model(sm).WherePKs(pks).Exec(ctx); sErr != nil {
		return sErr
	}
	if smRfl.ChainValue().Len() != len(pks) {
		panic("bad index length")
	}
	if tscErr := tsb.storage.NewTSCreate().Model(sm).Series().Exec(
		ctx); tscErr != nil {
		return tscErr
	}
	return nil
}
