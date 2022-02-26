package errutil

import "context"

// |||| CATCHER ||||

type Catcher struct {
	err error
}

type ActionFunc func() error

func (c *Catcher) Exec(actionFunc ActionFunc) {
	if c.err != nil {
		return
	}
	err := actionFunc()
	if err != nil {
		c.err = err
	}
}

func (c *Catcher) Reset() {
	c.err = nil
}

func (c *Catcher) Error() error {
	return c.err
}

type ContextCatcher struct {
	*Catcher
	ctx context.Context
}

func NewContextCatcher(ctx context.Context) *ContextCatcher {
	return &ContextCatcher{Catcher: &Catcher{}, ctx: ctx}
}

type ActionFuncContext func(ctx context.Context) error

func (c *ContextCatcher) Exec(actionFunc ActionFuncContext) {
	c.Catcher.Exec(func() error { return actionFunc(c.ctx) })
}
