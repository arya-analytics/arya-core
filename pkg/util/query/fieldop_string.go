// Code generated by "stringer -type=FieldFilter"; DO NOT EDIT.

package query

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FieldFilterGreaterThan-0]
	_ = x[FieldFilterLessThan-1]
	_ = x[FieldFilterInRange-2]
	_ = x[FilterFilterIsIn-3]
}

const _FieldOp_name = "FieldOpGreaterThanFieldOpLessThanFieldOpInRangeFieldOpIn"

var _FieldOp_index = [...]uint8{0, 18, 33, 47, 56}

func (i FieldFilter) String() string {
	if i < 0 || i >= FieldFilter(len(_FieldOp_index)-1) {
		return "FieldFilter(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldOp_name[_FieldOp_index[i]:_FieldOp_index[i+1]]
}
