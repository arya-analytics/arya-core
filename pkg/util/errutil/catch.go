package errutil

import (
	"context"
	"reflect"
)

type Catch interface {
	Error() error
	Errors() []error
	Reset()
	AddError(args ...interface{})
}

// |||| OPTS ||||

type catchOpts struct {
	aggregate bool
	convert   ConvertChain
	hooks     []CatchHook
}

type CatchOpt func(o *catchOpts)

func WithAggregation() CatchOpt {
	return func(o *catchOpts) {
		o.aggregate = true
	}
}

func WithConvert(cc ConvertChain) CatchOpt {
	return func(o *catchOpts) {
		o.convert = cc
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
		c.errors = append(c.errors, c.convert(err))
	}
}

func (c *CatchSimple) convert(err error) error {
	if c.opts.convert != nil {
		return c.opts.convert.Exec(err)
	}
	return err
}

func (c *CatchSimple) AddError(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	la := args[len(args)-1]
	if la == nil {
		return
	}
	err, ok := la.(error)
	if !ok {
		if reflect.TypeOf(args[0]).Kind() == reflect.Func {
			panic("function must be called when running AddError!")
		}
		panic("catch function didn't return an error as its last value!")
	}
	c.Exec(func() error { return err })
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

// |||| HOOKS ||||

// || PIPE ||

func NewPipeHook(pipe chan error) func(err error) {
	return func(err error) {
		pipe <- err
	}
}
