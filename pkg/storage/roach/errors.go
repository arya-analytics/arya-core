package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/pg"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

var _ErrorTypeConverterChain = storage.ErrorTypeConverterChain{
	errConverterPG,
}
var _defaultConverter = errConverterDefault

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(_ErrorTypeConverterChain, _defaultConverter)
}

var _sqlErrors = map[string]storage.ErrorType{
	"sql: no rows in result set":                  storage.ErrorTypeItemNotFound,
	"constraint failed: UNIQUE constraint failed": storage.ErrorTypeUniqueViolation,
	"SQL logic errutil: no such table":            storage.ErrorTypeMigration,
	"bun: Update and Delete queries require at least one Where": storage.
		ErrorTypeInvalidArgs,
}

func errConverterDefault(err error) (storage.ErrorType, bool) {
	for k, v := range _sqlErrors {
		if strings.Contains(err.Error(), k) {
			return v, true
		}
	}
	return storage.ErrorTypeUnknown, false
}

var _pgErrs = map[pg.ErrorType]storage.ErrorType{
	pg.ErrorTypeUniqueViolation:     storage.ErrorTypeUniqueViolation,
	pg.ErrorTypeForeignKeyViolation: storage.ErrorTypeRelationshipViolation,
	pg.ErrorTypeIntegrityConstraint: storage.ErrorTypeInvalidField,
}

func errConverterPG(err error) (storage.ErrorType, bool) {
	driverErr, ok := err.(pgdriver.Error)
	if !ok {
		return storage.ErrorTypeUnknown, false
	}
	pgErr := pg.NewError(driverErr.Field('C'))
	ot, ok := _pgErrs[pgErr.Type]
	return ot, ok
}
