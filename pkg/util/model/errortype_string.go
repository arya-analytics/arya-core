// Code generated by "stringer -type=ErrorType"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrTypeNonPointer-0]
	_ = x[ErrTypeNonStructOrSlice-1]
	_ = x[ErrTypeIncompatibleModels-2]
}

const _ErrorType_name = "ErrTypeNonPointerErrTypeNonStructOrSliceErrTypeIncompatibleModels"

var _ErrorType_index = [...]uint8{0, 17, 40, 65}

func (i ErrorType) String() string {
	if i < 0 || i >= ErrorType(len(_ErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorType_name[_ErrorType_index[i]:_ErrorType_index[i+1]]
}
