package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type runQuery struct {
	s *storage
}

func newRunQuery(s *storage) *runQuery {
	return &runQuery{s: s}
}

func (rq *runQuery) runBeforeHooks(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, hook := range rq.s.queryHooks {
		c.Exec(hook.BeforeQuery)
	}
	return c.Error()
}

func (rq *runQuery) runAfterHooks(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, hook := range rq.s.queryHooks {
		c.Exec(hook.AfterQuery)
	}
	return c.Error()
}

func (rq *runQuery) exec(ctx context.Context, p *query.Pack) error {
	qc := query.NewCatch(ctx, p)
	qc.Exec(rq.runBeforeHooks)
	qc.Exec(rq.s.cfg.EngineMD.Exec)
	qc.Exec(rq.s.cfg.EngineObject.Exec)
	qc.Exec(rq.s.cfg.EngineCache.Exec)
	qc.Exec(rq.runAfterHooks)
	return qc.Error()
}
