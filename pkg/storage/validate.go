package storage

import "reflect"

func validateSliceOrStruct(v reflect.Value) error {
	k := v.Type().Elem().Kind()
	if k != reflect.Struct && k != reflect.Slice {
		return NewError(ErrTypeNonStructOrSlice)
	}
	return nil
}

func validateContainerIsPointer(v reflect.Value) error {
	k := v.Type().Kind()
	if k != reflect.Ptr {
		return NewError(ErrTypeNonPointer)
	}
	return nil
}
