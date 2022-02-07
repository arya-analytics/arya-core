package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/pg"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

var _errTypeConverterChain = storage.ErrorTypeConverterChain{
	errConverterPG,
}
var _defaultConverter = errConverterDefault

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(_errTypeConverterChain, _defaultConverter)
}

var _sqlErrors = map[string]storage.ErrorType{
	"sql: no rows in result set":                  storage.ErrTypeItemNotFound,
	"constraint failed: UNIQUE constraint failed": storage.ErrTypeUniqueViolation,
	"SQL logic errutil: no such table":            storage.ErrTypeMigration,
	"bun: Update and Delete queries require at least one Where": storage.
		ErrTypeInvalidArgs,
}

func errConverterDefault(err error) (storage.ErrorType, bool) {
	for k, v := range _sqlErrors {
		if strings.Contains(err.Error(), k) {
			return v, true
		}
	}
	return storage.ErrTypeUnknown, false
}

var _pgErrs = map[pg.ErrorType]storage.ErrorType{
	pg.ErrTypeUniqueViolation:     storage.ErrTypeUniqueViolation,
	pg.ErrTypeForeignKeyViolation: storage.ErrTypeRelationshipViolation,
	pg.ErrTypeIntegrityConstraint: storage.ErrTypeInvalidField,
}

func errConverterPG(err error) (storage.ErrorType, bool) {
	driverErr, ok := err.(pgdriver.Error)
	if !ok {
		return storage.ErrTypeUnknown, false
	}
	pgErr := pg.NewError(driverErr.Field('C'))
	ot, ok := _pgErrs[pgErr.Type]
	return ot, ok
}
