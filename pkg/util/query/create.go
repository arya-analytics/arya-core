package query

// Create is used for writing Queries that create/persist models to storage.
type Create struct {
	Base
}

// || CONSTRUCTOR ||

// NewCreate instantiates a new Create query.
func NewCreate() *Create {
	c := &Create{}
	c.Base.Init(c)
	return c
}

// || MODEL ||

// Model sets the model to create. model must be passed as a pointer.
// The model can be a pointer to a struct or a pointer to a slice.
// The model must contain all necessary values and satisfy any relationships.
func (c *Create) Model(m interface{}) *Create {
	c.Base.Model(m)
	return c
}

// || EXEC ||

// BindExec binds Execute that Create will use to run the query.
// This method MUST be called before calling Exec.
func (c *Create) BindExec(e Execute) *Create {
	c.Base.BindExec(e)
	return c
}
