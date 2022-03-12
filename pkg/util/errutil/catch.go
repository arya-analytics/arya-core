package errutil

import (
	"context"
)

type Catch interface {
	Error() error
	Errors() []error
	AddHook(hook CatchHook)
	Reset()
}

// |||| OPTS ||||

type catchOpts struct {
	aggregate bool
	hooks     []CatchHook
}

type CatchOpt func(o *catchOpts)

func WithAggregation() CatchOpt {
	return func(o *catchOpts) {
		o.aggregate = true
	}
}

type CatchHook func(err error)

func WithHooks(hooks ...CatchHook) CatchOpt {
	return func(o *catchOpts) {
		o.hooks = hooks
	}
}

// |||| SIMPLE ||||

type CatchSimple struct {
	errors []error
	opts   *catchOpts
}

func NewCatchSimple(opts ...CatchOpt) *CatchSimple {
	c := &CatchSimple{opts: &catchOpts{}}
	c.bindOpts(opts...)
	return c
}

type CatchAction func() error

func (c *CatchSimple) bindOpts(opts ...CatchOpt) {
	for _, o := range opts {
		o(c.opts)
	}
}

func (c *CatchSimple) Exec(ca CatchAction) {
	if !c.opts.aggregate && len(c.errors) > 0 {
		return
	}
	err := ca()
	if err != nil {
		c.runHooks(err)
		c.errors = append(c.errors, err)
	}
}

func (c *CatchSimple) AddHook(hook CatchHook) {
	c.bindOpts(WithHooks(hook))
}

func (c *CatchSimple) runHooks(err error) {
	for _, h := range c.opts.hooks {
		h(err)
	}
}

func (c *CatchSimple) Reset() {
	c.errors = []error{}
}

func (c *CatchSimple) Error() error {
	if len(c.Errors()) == 0 {
		return nil
	}
	return c.Errors()[0]
}

func (c *CatchSimple) Errors() []error {
	return c.errors
}

// |||| CONTEXT ||||

type CatchContext struct {
	*CatchSimple
	ctx context.Context
}

func NewCatchContext(ctx context.Context, opts ...CatchOpt) *CatchContext {
	return &CatchContext{CatchSimple: NewCatchSimple(opts...), ctx: ctx}
}

type CatchActionCtx func(ctx context.Context) error

func (c *CatchContext) Exec(ca CatchActionCtx) {
	c.CatchSimple.Exec(func() error { return ca(c.ctx) })
}

// |||| PIPE ||||

func NewPipeHook(pipe chan error) func(err error) {
	return func(err error) {
		pipe <- err
	}
}
