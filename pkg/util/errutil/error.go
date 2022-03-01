package errutil

import "context"

type Catch interface {
	Exec(actionFunc CatchAction)
	Error() error
	Errors() []error
}

// |||| SIMPLE CATCH ||||

type CatchSimple struct {
	err error
}

type CatchAction func() error

func (c *CatchSimple) Exec(ca CatchAction) {
	if c.err != nil {
		return
	}
	err := ca()
	if err != nil {
		c.err = err
	}
}

func (c *CatchSimple) Reset() {
	c.err = nil
}

func (c *CatchSimple) Error() error {
	return c.err
}

func (c *CatchSimple) Errors() []error {
	return []error{c.err}
}

// |||| CATCH W CONTEXT ||||

type CatchWCtx struct {
	*CatchSimple
	ctx context.Context
}

func NewCatchWCtx(ctx context.Context) *CatchWCtx {
	return &CatchWCtx{CatchSimple: &CatchSimple{}, ctx: ctx}
}

type CatchActionCtx func(ctx context.Context) error

func (c *CatchWCtx) Exec(ca CatchActionCtx) {
	c.CatchSimple.Exec(func() error { return ca(c.ctx) })
}

// |||| CATCH AGGREGATE |||

type CatchAggregate struct {
	errors []error
}

func (c *CatchAggregate) Exec(actionFunc CatchAction) {
	err := actionFunc()
	if err != nil {
		c.errors = append(c.errors, err)
	}
}

func (c *CatchAggregate) Error() error {
	if len(c.Errors()) == 0 {
		return nil
	}
	return c.Errors()[0]
}

func (c *CatchAggregate) Errors() []error {
	return c.errors
}
