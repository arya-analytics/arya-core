package query

import (
	"context"
	"fmt"
)

// Execute represents a function or method that can execute the provided Pack
// on a persistent data store. It should return an error representing any errors encountered
// during query execution.
//
// If you want to call different Execute based on the type of query, see Switch.
type Execute func(ctx context.Context, p *Pack) error

// |||| SWITCH ||||

// Ops represents a set of Execute to call for a specific Query.
type Ops struct {
	Create   Execute
	Retrieve Execute
	Delete   Execute
	Update   Execute
}

// Switch switches a Pack to a configured set of Execute. Switch allows the caller to implement different query Execute
// depending on the Query used.
func Switch(ctx context.Context, p *Pack, ops Ops) error {
	switch p.Query().(type) {
	case *Create:
		if ops.Create != nil {
			return ops.Create(ctx, p)
		}
	case *Retrieve:
		if ops.Retrieve != nil {
			return ops.Retrieve(ctx, p)
		}
	case *Update:
		if ops.Update != nil {
			return ops.Update(ctx, p)
		}
	case *Delete:
		if ops.Delete != nil {
			return ops.Delete(ctx, p)
		}
	}
	panic(fmt.Sprintf("%T not supported for model %s", p.Query(), p.Model().Type().Name()))
}
