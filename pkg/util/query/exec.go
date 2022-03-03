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
//
// Parsing and Executing a packed Query (Pack).
//
// 1. To parse the Pack, use the different opts available to retrieve options from the Pack. As an example, here's how
// to retrieve the primary key of a query.
//
//		pkc, ok := query.PKOpt(p)
//
// pkc will hold a model.PKChain representing the primary keys of the query. ok will be false if the primary keys don't
// exist. Repeat this process with the different options you want to provide support for to extract all the info you need.
// The different options available for queries are suffixed with the 'Opt' keyword.
//
// 2. Use the parsed options and provided context to run the query, and binds the results into the Pack.Model().
// in the case of options with an 'into' arg (like CalcOpt), bind the result into the provided argument.
//
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
