package query

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	log "github.com/sirupsen/logrus"
)

const (
	errKey = "query"
)

type Error struct {
	Base    error
	Type    ErrorType
	Message string
}

func NewSimpleError(t ErrorType, base error) Error {
	e := Error{Type: t, Base: base}
	if e.Base != nil {
		e.Message = e.Base.Error()
	}
	return e
}

func NewUnknownError(base error) Error {
	return NewSimpleError(ErrorTypeUnknown, base)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s - %s - %s", errKey, e.Type, e.Message, e.Base)
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeItemNotFound
	ErrorTypeUniqueViolation
	ErrorTypeRelationshipViolation
	ErrorTypeInvalidField
	ErrorTypeMigration
	ErrorTypeInvalidArgs
	ErrorTypeConnection
	ErrorTypeMultipleResults
)

func injectErrKey(errStr string, args ...interface{}) string {
	return fmt.Sprintf("%s -> %s", errKey, fmt.Sprintf(errStr, args...))
}

// |||| CONVERTER ||||

// NewErrorConvertChain wraps errutil.ConvertChain and adds the following errutil.Convert
// implementations:
//
// 1. A pass through errutil.Convert that will propagate the error if it is already of type query.Error.
//
// 2. General errutil.Convert that handle common query errors
//
// 3. A default errutil.Convert that will return a query.Error with query.ErrorTypeUnknown.
//
func NewErrorConvertChain(converters ...errutil.Convert) errutil.ConvertChain {
	cc := errutil.ConvertChain{errorPassConvert}
	cc = append(cc, converters...)
	cc = append(cc, errorContextCanceled, errorDefaultConvert)
	return cc
}

func errorPassConvert(err error) (error, bool) {
	_, ok := err.(Error)
	return err, ok
}

func errorContextCanceled(err error) (error, bool) {
	if err.Error() == "context canceled" {
		return NewSimpleError(ErrorTypeInvalidArgs, err), true
	}
	return err, false
}

func errorDefaultConvert(err error) (error, bool) {
	log.Errorf(injectErrKey("unknown error -> %s", err))
	return NewUnknownError(err), true
}

// |||| CATCH ||||

// Catch wraps errutil.CatchContext to help running contiguous sets of Execute (i.e. executing multiple Query in a row)
// Catch supplements errutil.CatchContext context by providing a Pack as well.
type Catch struct {
	p *Pack
	*errutil.CatchContext
}

// NewCatch creates a new catch with the provided context.Context and Pack.
func NewCatch(ctx context.Context, p *Pack, opts ...errutil.CatchOpt) *Catch {
	return &Catch{CatchContext: errutil.NewCatchContext(ctx, opts...), p: p}
}

// Exec runs the provided Execute and catches an of the errors returned.
func (c *Catch) Exec(exec Execute) {
	c.CatchContext.Exec(func(ctx context.Context) error { return exec(ctx, c.p) })
}
