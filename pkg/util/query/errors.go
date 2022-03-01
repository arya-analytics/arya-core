package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Catch struct {
	p *Pack
	*errutil.CatchWCtx
}

type CatchAction func(ctx context.Context, p *Pack) error

func NewCatch(ctx context.Context, p *Pack) *Catch {
	return &Catch{CatchWCtx: errutil.NewCatchWCtx(ctx), p: p}
}

func (c *Catch) Exec(ca CatchAction) {
	c.CatchWCtx.Exec(func(ctx context.Context) error { return ca(ctx, c.p) })
}
