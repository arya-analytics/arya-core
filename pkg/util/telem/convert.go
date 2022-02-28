package telem

import (
	"encoding/binary"
	"math"
	"reflect"
)

func ByteOrder() binary.ByteOrder {
	return binary.BigEndian
}

func convertBytes(b []byte, dataType DataType) interface{} {
	ss := int(sampleSize(dataType))
	rv := newSliceOfDataType(dataType, len(b))
	cv := convertCatalog()[dataType]
	for i := 0; i < len(b); i += ss {
		rvI := i / ss
		rv.Index(rvI).Set(reflect.ValueOf(cv(b[i : i+ss])))
	}
	if len(b) == ss {
		return rv.Index(0).Interface()
	}
	return rv.Interface()
}

/// |||| CATALOG ||||

func newSliceOfDataType(dataType DataType, len int) reflect.Value {
	ss := int(sampleSize(dataType))
	rv := reflect.MakeSlice(dataTypeCatalog()[dataType], len/ss, len)
	rvPtr := reflect.New(rv.Type())
	rvPtr.Elem().Set(rv)
	return rv
}

func dataTypeCatalog() map[DataType]reflect.Type {
	return map[DataType]reflect.Type{
		DataTypeFloat32: reflect.TypeOf([]float32{}),
		DataTypeFloat64: reflect.TypeOf([]float64{}),
	}
}

/// |||| SINGLE BYTE CONVERTERS ||||

type convert func(b []byte) interface{}

func convertCatalog() map[DataType]convert {
	return map[DataType]convert{
		DataTypeFloat32: sbToFloat32,
		DataTypeFloat64: sbToFloat64,
	}
}

func sbToFloat32(b []byte) interface{} {
	return math.Float32frombits(ByteOrder().Uint32(b))
}

func sbToFloat64(b []byte) interface{} {
	return math.Float64frombits(ByteOrder().Uint64(b))
}
