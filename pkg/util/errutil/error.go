package errutil

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
