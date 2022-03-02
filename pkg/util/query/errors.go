package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

// Catch wraps errutil.CatchWCtx to help running contiguous sets of Execute (i.e. executing multiple Query in a row)
// Catch supplements errutil.CatchWCtx context by providing a Pack as well.
type Catch struct {
	p *Pack
	*errutil.CatchWCtx
}

// NewCatch creates a new catch with the provided context.Context and Pack.
func NewCatch(ctx context.Context, p *Pack) *Catch {
	return &Catch{CatchWCtx: errutil.NewCatchWCtx(ctx), p: p}
}

// Exec runs the provided Execute and catches an of the errors returned.
func (c *Catch) Exec(exec Execute) {
	c.CatchWCtx.Exec(func(ctx context.Context) error { return exec(ctx, c.p) })
}
