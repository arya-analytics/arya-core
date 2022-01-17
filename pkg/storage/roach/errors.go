package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/pg"
	"github.com/uptrace/bun/driver/pgdriver"
)

// |||| BUN ERRORS ||||

func parseBunErr(err error) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case pgdriver.Error:
		return parsePgDriverErr(err.(pgdriver.Error))
	default:
		return parseSqlError(err.Error())
	}
}

// |||| PGDRIVER ERRORS ||||

func parsePgDriverErr(err pgdriver.Error) error {
	pgErr := pg.NewError(err.Field('C'))
	return storage.NewError(pgToStorageErrType(pgErr.Type))
}

var _pgErrs = map[pg.ErrorType]storage.ErrorType{
	pg.ErrTypeUniqueViolation:     storage.ErrTypeUniqueViolation,
	pg.ErrTypeForeignKeyViolation: storage.ErrTypeRelationshipViolation,
	pg.ErrTypeIntegrityConstraint: storage.ErrTypeInvalidField,
}

func pgToStorageErrType(t pg.ErrorType) storage.ErrorType {
	ot, ok := _pgErrs[t]
	if !ok {
		return storage.ErrTypeUnknown
	}
	return ot
}

// |||| SQL ERRORS |||

var _sqlErrors = map[string]storage.ErrorType{
	"sql: no rows in result set": storage.ErrTypeItemNotFound,
}

func sqlToStorageErr(sql string) storage.ErrorType {
	ot, ok := _sqlErrors[sql]
	if !ok {
		return storage.ErrTypeUnknown
	}
	return ot
}

func parseSqlError(sql string) storage.Error {
	return storage.NewError(sqlToStorageErr(sql))
}
