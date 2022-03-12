package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type def struct {
	s *storage
}

func newDef(s *storage) *def {
	return &def{s: s}
}

type update struct {
	def
}

func newUpdate(s *storage) *update {
	return &update{def: def{s: s}}
}

func (d *def) runBeforeHooks(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, hook := range d.s.queryHooks {
		c.Exec(hook.BeforeQuery)
	}
	return c.Error()
}

func (d *def) runAfterHooks(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p, errutil.WithAggregation())
	for _, hook := range d.s.queryHooks {
		c.Exec(hook.AfterQuery)
	}
	return c.Error()
}

func (d *def) exec(ctx context.Context, p *query.Pack) error {
	qc := query.NewCatch(ctx, p)
	qc.Exec(d.runBeforeHooks)
	qc.Exec(d.s.cfg.EngineMD.Exec)
	qc.Exec(d.s.cfg.EngineObject.Exec)
	qc.Exec(d.runAfterHooks)
	return qc.Error()
}

func (u *update) exec(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p)
	c.Exec(u.runBeforeHooks)
	c.Exec(u.s.cfg.EngineMD.Exec)
	c.Exec(u.runAfterHooks)
	return c.Error()
}
