// Code generated by "stringer -type=TimingErrorType"; DO NOT EDIT.

package chanchunk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TimingErrorTypeChunkOverlap-1]
	_ = x[TimingErrorTypeIncompatibleChunks-2]
	_ = x[TimingErrorTypeNonContiguous-3]
}

const _TimingErrorType_name = "TimingErrorTypeChunkOverlapTimingErrorTypeIncompatibleChunksTimingErrorTypeNonContiguous"

var _TimingErrorType_index = [...]uint8{0, 27, 60, 88}

func (i TimingErrorType) String() string {
	i -= 1
	if i < 0 || i >= TimingErrorType(len(_TimingErrorType_index)-1) {
		return "TimingErrorType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TimingErrorType_name[_TimingErrorType_index[i]:_TimingErrorType_index[i+1]]
}
