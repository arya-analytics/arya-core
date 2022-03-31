package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Hook interface {
	Before(ctx context.Context, p *Pack) error
	After(ctx context.Context, p *Pack) error
}

type HookRunner struct {
	hooks map[Hook]bool
}

func NewHookRunner() *HookRunner {
	return &HookRunner{hooks: make(map[Hook]bool)}
}

func (hr *HookRunner) AddQueryHook(hook Hook) {
	hr.hooks[hook] = true
}

func (hr *HookRunner) RemoveQueryHook(hook Hook) {
	delete(hr.hooks, hook)
}

func (hr *HookRunner) Before(ctx context.Context, p *Pack) error {
	c := NewCatch(ctx, p, errutil.WithAggregation())
	for h := range hr.hooks {
		c.Exec(h.Before)
	}
	return c.Error()
}

func (hr *HookRunner) After(ctx context.Context, p *Pack) error {
	c := NewCatch(ctx, p, errutil.WithAggregation())
	for h := range hr.hooks {
		c.Exec(h.After)
	}
	return c.Error()
}

func (hr *HookRunner) ClearQueryHooks() {
	hr.hooks = make(map[Hook]bool)
}
