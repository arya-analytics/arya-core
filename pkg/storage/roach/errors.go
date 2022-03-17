package roach

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/pg"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/lib/pq"
	"github.com/uptrace/bun/driver/pgdriver"
	"net"
	"strings"
)

func newErrorConvert() errutil.ConvertChain {
	return query.NewErrorConvertChain(
		errorConvertPQ,
		errorConvertPGDriver,
		errorConvertConnection,
		errorConvertDefault,
	)
}

func errorConvertDefault(err error) (error, bool) {
	for k, t := range sqlErrors() {
		if strings.Contains(err.Error(), k) {
			return query.NewSimpleError(t, err), true
		}
	}
	return query.NewUnknownError(err), false
}

func sqlErrors() map[string]query.ErrorType {
	return map[string]query.ErrorType{
		"sql: no rows in result set":                  query.ErrorTypeItemNotFound,
		"constraint failed: UNIQUE constraint failed": query.ErrorTypeUniqueViolation,
		"SQL logic errutil: no such table":            query.ErrorTypeMigration,
		"bun: Update and Delete queries require at":   query.ErrorTypeInvalidArgs,
		"does not have relation":                      query.ErrorTypeInvalidArgs,
	}
}

func errorConvertConnection(err error) (error, bool) {
	switch err.(type) {
	case *net.OpError:
		return query.NewSimpleError(query.ErrorTypeConnection, err), true
	default:
		return err, false
	}

}

func pgErrors() map[pg.ErrorType]query.ErrorType {
	return map[pg.ErrorType]query.ErrorType{
		pg.ErrorTypeUniqueViolation:     query.ErrorTypeUniqueViolation,
		pg.ErrorTypeForeignKeyViolation: query.ErrorTypeRelationshipViolation,
		pg.ErrorTypeIntegrityConstraint: query.ErrorTypeInvalidField,
		pg.ErrorTypeSyntax:              query.ErrorTypeInvalidArgs,
	}
}

const (
	pgDriverCodeField = 'C'
)

func errorConvertPGDriver(err error) (error, bool) {
	driverErr, ok := err.(pgdriver.Error)
	if !ok {
		return query.NewUnknownError(err), false
	}
	pgErr := pg.NewError(driverErr.Field(pgDriverCodeField))
	t, ok := pgErrors()[pgErr.Type]
	return query.NewSimpleError(t, err), ok
}

func errorConvertPQ(err error) (error, bool) {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return query.NewUnknownError(err), false
	}
	pgErr := pg.NewError(string(pqErr.Code))
	t, ok := pgErrors()[pgErr.Type]
	return query.NewSimpleError(t, err), ok
}
