package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

// Hook is a callback interface that can be used to introspect or modify the parameters or results of a query.
type Hook interface {
	// Before is invoked before the query is executed.
	Before(ctx context.Context, p *Pack) error
	// After is invoked after the query is executed.
	After(ctx context.Context, p *Pack) error
}

// HookRunner wraps a set of hooks and provides functionality to execute them in sequence.
type HookRunner struct {
	hooks map[Hook]bool
}

// NewHookRunner instantiates a new empty HookRunner.
func NewHookRunner() *HookRunner {
	return &HookRunner{hooks: make(map[Hook]bool)}
}

// AddQueryHook adds a new Hook to the HookRunner.
func (hr *HookRunner) AddQueryHook(hook Hook) {
	hr.hooks[hook] = true
}

// RemoveQueryHook removes a Hook from the HookRunner.
func (hr *HookRunner) RemoveQueryHook(hook Hook) {
	delete(hr.hooks, hook)
}

// Before runs all Hook.Before hooks in sequence.
func (hr *HookRunner) Before(ctx context.Context, p *Pack) error {
	c := NewCatch(ctx, p, errutil.WithAggregation())
	for h := range hr.hooks {
		c.Exec(h.Before)
	}
	return c.Error()
}

// After runs all Hook.After hooks in sequence.
func (hr *HookRunner) After(ctx context.Context, p *Pack) error {
	c := NewCatch(ctx, p, errutil.WithAggregation())
	for h := range hr.hooks {
		c.Exec(h.After)
	}
	return c.Error()
}

// ClearQueryHooks removes all Hook from the HookRunner.
func (hr *HookRunner) ClearQueryHooks() {
	hr.hooks = make(map[Hook]bool)
}
