// Code generated by "stringer -type=FieldExpOp"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FieldExpOpGreaterThan-0]
	_ = x[FieldExpOpLessThan-1]
	_ = x[FieldExpOpInRange-2]
}

const _FieldExpOp_name = "FieldExpOpGreaterThanFieldExpOpLessThanFieldExpOpInRange"

var _FieldExpOp_index = [...]uint8{0, 21, 39, 56}

func (i FieldExpOp) String() string {
	if i < 0 || i >= FieldExpOp(len(_FieldExpOp_index)-1) {
		return "FieldExpOp(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldExpOp_name[_FieldExpOp_index[i]:_FieldExpOp_index[i+1]]
}
