package errutil

import "context"

type Catch interface {
	Exec(actionFunc ActionFunc)
	Error() error
	Errors() []error
}

// |||| SIMPLE CATCH ||||

type CatchSimple struct {
	err error
}

type ActionFunc func() error

func (c *CatchSimple) Exec(actionFunc ActionFunc) {
	if c.err != nil {
		return
	}
	err := actionFunc()
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

type CatchWContext struct {
	*CatchSimple
	ctx context.Context
}

func NewCatchWContext(ctx context.Context) *CatchWContext {
	return &CatchWContext{CatchSimple: &CatchSimple{}, ctx: ctx}
}

type ActionFuncContext func(ctx context.Context) error

func (c *CatchWContext) Exec(actionFunc ActionFuncContext) {
	c.CatchSimple.Exec(func() error { return actionFunc(c.ctx) })
}

// |||| CATCH AGGREGATE |||

type CatchAggregate struct {
	errors []error
}

func (c *CatchAggregate) Exec(actionFunc ActionFunc) {
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
