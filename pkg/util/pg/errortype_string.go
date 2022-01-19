// Code generated by "stringer -type=ErrorType"; DO NOT EDIT.

package pg

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrTypeUnknown-0]
	_ = x[ErrTypeWarning-1]
	_ = x[ErrTypeNoData-2]
	_ = x[ErrTypeSQLNotComplete-3]
	_ = x[ErrTypeConn-4]
	_ = x[ErrTypeTriggeredAction-5]
	_ = x[ErrTypeInvalidTransaction-6]
	_ = x[ErrTypeLocator-7]
	_ = x[ErrTypeGrantor-8]
	_ = x[ErrTypeInvalidRoleSpec-9]
	_ = x[ErrTypeDiagnosticsException-10]
	_ = x[ErrTypeCaseNotFound-11]
	_ = x[ErrTypeCardinalityViolation-12]
	_ = x[ErrTypeDataException-13]
	_ = x[ErrTypeIntegrityConstraint-14]
	_ = x[ErrTypeInvalidCursor-15]
	_ = x[ErrTypeInvalidSQLStatementName-16]
	_ = x[ErrTypeTriggeredDataChangeViolation-17]
	_ = x[ErrTypeInvalidAuthSpec-18]
	_ = x[ErrTypeDependentPrivileged-19]
	_ = x[ErrTypeSQLRoutine-20]
	_ = x[ErrTypeInvalidCatalog-21]
	_ = x[ErrTypeTransactionRollback-22]
	_ = x[ErrTypeSyntax-23]
	_ = x[ErrTypeSystem-24]
	_ = x[ErrTypeFeatureNotSupported-25]
	_ = x[ErrTypeUniqueViolation-26]
	_ = x[ErrTypeForeignKeyViolation-27]
}

const _ErrorType_name = "ErrTypeUnknownErrTypeWarningErrTypeNoDataErrTypeSQLNotCompleteErrTypeConnErrTypeTriggeredActionErrTypeInvalidTransactionErrTypeLocatorErrTypeGrantorErrTypeInvalidRoleSpecErrTypeDiagnosticsExceptionErrTypeCaseNotFoundErrTypeCardinalityViolationErrTypeDataExceptionErrTypeIntegrityConstraintErrTypeInvalidCursorErrTypeInvalidSQLStatementNameErrTypeTriggeredDataChangeViolationErrTypeInvalidAuthSpecErrTypeDependentPrivilegedErrTypeSQLRoutineErrTypeInvalidCatalogErrTypeTransactionRollbackErrTypeSyntaxErrTypeSystemErrTypeFeatureNotSupportedErrTypeUniqueViolationErrTypeForeignKeyViolation"

var _ErrorType_index = [...]uint16{0, 14, 28, 41, 62, 73, 95, 120, 134, 148, 170, 197, 216, 243, 263, 289, 309, 339, 374, 396, 422, 439, 460, 486, 499, 512, 538, 560, 586}

func (i ErrorType) String() string {
	if i < 0 || i >= ErrorType(len(_ErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorType_name[_ErrorType_index[i]:_ErrorType_index[i+1]]
}
