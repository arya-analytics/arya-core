// Code generated by "stringer -type=ErrorType"; DO NOT EDIT.

package api

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrorTypeUnknown-0]
	_ = x[ErrorTypeUnauthorized-1]
	_ = x[ErrorTypeAuthentication-2]
	_ = x[ErrorTypeInvalidArguments-3]
}

const _ErrorType_name = "ErrorTypeUnknownErrorTypeUnauthorizedErrorTypeAuthenticationErrorTypeInvalidArguments"

var _ErrorType_index = [...]uint8{0, 16, 37, 60, 85}

func (i ErrorType) String() string {
	if i < 0 || i >= ErrorType(len(_ErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorType_name[_ErrorType_index[i]:_ErrorType_index[i+1]]
}
