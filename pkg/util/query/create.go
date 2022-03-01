package query

// QueryCreate creates a new model in storage.
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

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct or a pointer to a slice.
// The model must contain all necessary values and satisfy any relationships.
func (c *Create) Model(m interface{}) *Create {
	c.baseModel(m)
	return c
}

// || EXEC ||

func (c *Create) BindExec(e Execute) *Create {
	c.baseBindExec(e)
	return c
}
