// Code generated by "stringer -type=ErrorType"; DO NOT EDIT.

package chanchunk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrorTimingOverlap-1]
	_ = x[ErrorTimingIncompatibleChunks-2]
	_ = x[ErrorTimingNonContiguous-3]
}

const _TimingErrorType_name = "TimingErrorTypeChunkOverlapTimingErrorTypeIncompatibleChunksTimingErrorTypeNonContiguous"

var _TimingErrorType_index = [...]uint8{0, 27, 60, 88}

func (i ErrorType) String() string {
	i -= 1
	if i < 0 || i >= ErrorType(len(_TimingErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TimingErrorType_name[_TimingErrorType_index[i]:_TimingErrorType_index[i+1]]
}
