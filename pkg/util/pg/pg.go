package pg

import "fmt"

const (
	errKey         = "pg"
	errClassLength = 2
)

var _errs = map[string]ErrorType{
	"23505": ErrTypeUniqueViolation,
	"23503": ErrTypeForeignKeyViolation,
}

var _errClasses = map[string]ErrorType{
	"01": ErrTypeWarning,
	"02": ErrTypeNoData,
	"03": ErrTypeSQLNotComplete,
	"08": ErrTypeConn,
	"09": ErrTypeTriggeredAction,
	"0A": ErrTypeFeatureNotSupported,
	"0B": ErrTypeInvalidTransaction,
	"0F": ErrTypeLocator,
	"0L": ErrTypeGrantor,
	"OP": ErrTypeInvalidRoleSpec,
	"0Z": ErrTypeDiagnosticsException,
	"20": ErrTypeCaseNotFound,
	"21": ErrTypeCardinalityViolation,
	"22": ErrTypeDataException,
	"23": ErrTypeIntegrityConstraint,
	"24": ErrTypeInvalidCursor,
	"25": ErrTypeInvalidTransaction,
	"26": ErrTypeInvalidSQLStatementName,
	"27": ErrTypeTriggeredDataChangeViolation,
	"28": ErrTypeInvalidAuthSpec,
	"2B": ErrTypeDependentPrivileged,
	"2D": ErrTypeInvalidSQLStatementName,
	"2F": ErrTypeSQLRoutine,
	"34": ErrTypeInvalidCursor,
	"3D": ErrTypeInvalidCatalog,
	"40": ErrTypeTransactionRollback,
	"42": ErrTypeSyntax,
	"58": ErrTypeSystem,
}

func errTypeFromCode(code string) ErrorType {
	t, ok := _errs[code]
	if !ok {
		t, ok = _errClasses[code[0:errClassLength]]
		if !ok {
			t = ErrTypeUnknown
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
		Type: errTypeFromCode(code),
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", errKey, e.Type)
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeWarning
	ErrTypeNoData
	ErrTypeSQLNotComplete
	ErrTypeConn
	ErrTypeTriggeredAction
	ErrTypeInvalidTransaction
	ErrTypeLocator
	ErrTypeGrantor
	ErrTypeInvalidRoleSpec
	ErrTypeDiagnosticsException
	ErrTypeCaseNotFound
	ErrTypeCardinalityViolation
	ErrTypeDataException
	ErrTypeIntegrityConstraint
	ErrTypeInvalidCursor
	ErrTypeInvalidSQLStatementName
	ErrTypeTriggeredDataChangeViolation
	ErrTypeInvalidAuthSpec
	ErrTypeDependentPrivileged
	ErrTypeSQLRoutine
	ErrTypeInvalidCatalog
	ErrTypeTransactionRollback
	ErrTypeSyntax
	ErrTypeSystem
	ErrTypeFeatureNotSupported
	ErrTypeUniqueViolation
	ErrTypeForeignKeyViolation
)
