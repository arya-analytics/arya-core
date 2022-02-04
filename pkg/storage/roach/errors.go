package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/pg"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

// |||| BUN ERRORS ||||

func parseBunErr(err error) (oErr error) {
	if err == nil {
		return oErr
	}

	switch err := err.(type) {
	case pgdriver.Error:
		oErr = parsePgDriverErr(err)
	default:
		oErr = parseSqlError(err.Error())
	}
	se, ok := oErr.(storage.Error)
	if ok {
		if se.Type == storage.ErrTypeUnknown {
			log.Errorf("Unknown err -> %s", err)
		}
	}
	return oErr
}

// |||| PGDRIVER ERRORS ||||

func parsePgDriverErr(err pgdriver.Error) error {
	pgErr := pg.NewError(err.Field('C'))
	return storage.Error{Base: err, Type: pgToStorageErrType(pgErr.Type),
		Message: err.Error()}
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
	"sql: no rows in result set":                  storage.ErrTypeItemNotFound,
	"constraint failed: UNIQUE constraint failed": storage.ErrTypeUniqueViolation,
	"SQL logic errutil: no such table":            storage.ErrTypeMigration,
	"bun: Update and Delete queries require at least one Where": storage.
		ErrTypeInvalidArgs,
}

func sqlToStorageErr(sql string) storage.ErrorType {
	for k, v := range _sqlErrors {
		if strings.Contains(sql, k) {
			return v
		}
	}
	return storage.ErrTypeUnknown
}

func parseSqlError(sql string) error {
	return storage.Error{Type: sqlToStorageErr(sql), Message: sql}
}
