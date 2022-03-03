package errutil

import "context"

type Catch interface {
	Exec(actionFunc CatchAction)
	Error() error
	Errors() []error
	Reset()
}

// |||| OPTS ||||

type catchOpts struct {
	aggregate bool
}

type CatchOpt func(o *catchOpts)

func WithAggregation() CatchOpt {
	return func(o *catchOpts) {
		o.aggregate = true
	}
}

// |||| SIMPLE CATCH ||||

type CatchSimple struct {
	errors []error
	opts   *catchOpts
}

func NewCatchSimple(opts ...CatchOpt) *CatchSimple {
	c := &CatchSimple{opts: &catchOpts{}}
	for _, o := range opts {
		o(c.opts)
	}
	return c
}

type CatchAction func() error

func (c *CatchSimple) Exec(ca CatchAction) {
	if !c.opts.aggregate && len(c.errors) > 0 {
		return
	}
	err := ca()
	if err != nil {
		c.errors = append(c.errors, err)
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

// |||| CATCH W CONTEXT ||||

type CatchWCtx struct {
	*CatchSimple
	ctx context.Context
}

func NewCatchWCtx(ctx context.Context, opts ...CatchOpt) *CatchWCtx {
	return &CatchWCtx{CatchSimple: NewCatchSimple(opts...), ctx: ctx}
}

type CatchActionCtx func(ctx context.Context) error

func (c *CatchWCtx) Exec(ca CatchActionCtx) {
	c.CatchSimple.Exec(func() error { return ca(c.ctx) })
}
