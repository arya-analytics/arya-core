package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type QueryHook interface {
	Before(ctx context.Context, p *query.Pack) error
	After(ctx context.Context, p *query.Pack) error
}

type queryHookChain []QueryHook

func (qhc queryHookChain) before(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, h := range qhc {
		c.Exec(h.After)
	}
	return c.Error()
}

func (qhc queryHookChain) after(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, h := range qhc {
		c.Exec(h.Before)
	}
	return c.Error()
}
