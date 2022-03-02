package query

// Create is used for writing Queries that create/persist models to storage.
type Create struct {
	base
}

// || CONSTRUCTOR ||

// NewCreate instantiates a new Create query.
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
