package query

type Create struct {
	base
}

// || CONSTRUCTOR ||

func NewCreate() *Create {
	c := &Create{}
	c.baseInit(c)
	return c
}

// || MODEL ||

func (c *Create) Model(m interface{}) *Create {
	c.baseModel(m)
	return c
}

// || EXEC ||

func (c *Create) BindExec(e Execute) *Create {
	c.baseBindExec(e)
	return c
}
