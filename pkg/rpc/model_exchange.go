package rpc

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"reflect"
)

func NewModelExchange(sourcePtr, destPtr interface{}) *model.Exchange {
	return model.NewExchange(
		sourcePtr,
		destPtr,
		model.FieldHandlerPK,
		FieldHandlerTelemBulk,
	)
}

func FieldHandlerTelemBulk(sourceST, destST model.StructTag, sourceFld, destFld reflect.Value) (reflect.Value, bool) {
	if !taggedTelemBulk(sourceST) || !taggedTelemBulk(destST) {
		return reflect.Value{}, false
	}
	if isTelemBulkField(sourceFld) && isBytesField(destFld) {
		blk, b := sourceFld.Interface().(*telem.ChunkData), destFld.Interface().([]byte)
		if sourceFld.IsNil() {
			blk = telem.NewChunkData([]byte{})
			sourceFld.Set(reflect.ValueOf(blk))
		}
		if _, err := blk.Write(b); err != nil {
			panic(err)
		}
	} else if isBytesField(sourceFld) && isTelemBulkField(destFld) {
		blk := destFld.Interface().(*telem.ChunkData)
		sourceFld.Set(reflect.ValueOf(blk.Bytes()))
	} else {
		panic(fmt.Sprintf(
			"fields tagged telemChunkData, but didn't receive correct types! received %s and %s",
			sourceFld.Type(),
			destFld.Type(),
		))
	}
	return sourceFld, true
}

func isBytesField(fld reflect.Value) bool {
	return fld.Type() == reflect.TypeOf([]byte{}) || fld.Type() == reflect.TypeOf([]uint8{})
}

func isTelemBulkField(fld reflect.Value) bool {
	return fld.Type() == reflect.TypeOf(&telem.ChunkData{})
}

func taggedTelemBulk(st model.StructTag) bool {
	return st.Match(model.TagCat, model.RoleKey, "telemChunkData")
}
