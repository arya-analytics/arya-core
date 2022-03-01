package pg

import "fmt"

const (
	errKey         = "pg"
	errClassLength = 2
)

var _errs = map[string]ErrorType{
	"23505": ErrorTypeUniqueViolation,
	"23503": ErrorTypeForeignKeyViolation,
}

var _errClasses = map[string]ErrorType{
	"01": ErrorTypeWarning,
	"02": ErrorTypeNoData,
	"03": ErrorTypeSQLNotComplete,
	"08": ErrorTypeConn,
	"09": ErrorTypeTriggeredAction,
	"0A": ErrorTypeFeatureNotSupported,
	"0B": ErrorTypeInvalidTransaction,
	"0F": ErrorTypeLocator,
	"0L": ErrorTypeGrantor,
	"Op": ErrorTypeInvalidRoleSpec,
	"0Z": ErrorTypeDiagnosticsException,
	"20": ErrorTypeCaseNotFound,
	"21": ErrorTypeCardinalityViolation,
	"22": ErrorTypeDataException,
	"23": ErrorTypeIntegrityConstraint,
	"24": ErrorTypeInvalidCursor,
	"25": ErrorTypeInvalidTransaction,
	"26": ErrorTypeInvalidSQLStatementName,
	"27": ErrorTypeTriggeredDataChangeViolation,
	"28": ErrorTypeInvalidAuthSpec,
	"2B": ErrorTypeDependentPrivileged,
	"2D": ErrorTypeInvalidSQLStatementName,
	"2F": ErrorTypeSQLRoutine,
	"34": ErrorTypeInvalidCursor,
	"3D": ErrorTypeInvalidCatalog,
	"40": ErrorTypeTransactionRollback,
	"42": ErrorTypeSyntax,
	"58": ErrorTypeSystem,
}

func ErrorTypeFromCode(code string) ErrorType {
	t, ok := _errs[code]
	if !ok {
		t, ok = _errClasses[code[0:errClassLength]]
		if !ok {
			t = ErrorTypeUnknown
		}
	}
	return t
}

type Error struct {
	Code string
	Type ErrorType
}

func NewError(code string) Error {
	return Error{
		Code: code,
		Type: ErrorTypeFromCode(code),
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", errKey, e.Type)
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeWarning
	ErrorTypeNoData
	ErrorTypeSQLNotComplete
	ErrorTypeConn
	ErrorTypeTriggeredAction
	ErrorTypeInvalidTransaction
	ErrorTypeLocator
	ErrorTypeGrantor
	ErrorTypeInvalidRoleSpec
	ErrorTypeDiagnosticsException
	ErrorTypeCaseNotFound
	ErrorTypeCardinalityViolation
	ErrorTypeDataException
	ErrorTypeIntegrityConstraint
	ErrorTypeInvalidCursor
	ErrorTypeInvalidSQLStatementName
	ErrorTypeTriggeredDataChangeViolation
	ErrorTypeInvalidAuthSpec
	ErrorTypeDependentPrivileged
	ErrorTypeSQLRoutine
	ErrorTypeInvalidCatalog
	ErrorTypeTransactionRollback
	ErrorTypeSyntax
	ErrorTypeSystem
	ErrorTypeFeatureNotSupported
	ErrorTypeUniqueViolation
	ErrorTypeForeignKeyViolation
)
