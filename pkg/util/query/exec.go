package query

import (
	"context"
	"fmt"
	"reflect"
)

// Execute represents a function or method that can execute the provided Pack
// on a persistent data store. It should return an error representing any errors encountered
// during query execution.
//
// If you want to call different Execute based on the type of query, see Switch.
//
// Parsing and Executing a packed Query (Pack).
//
// 1. To parse the Pack, use the different Opts available to retrieve options from the Pack. As an example, here's how
// to retrieve the primary key of a query.
//
//		pkc, ok := query.RetrievePKOpt(p)
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

// Ops represents a set of Execute to call for a specific Query variant (Create, Retrieve, Update, etc).
type Ops map[Query]Execute

// Switch switches a Pack to a configured set of Execute. Switch allows the caller to implement different query Execute
// depending on the Query type used.
func Switch(ctx context.Context, p *Pack, ops Ops, opts ...SwitchOpt) error {
	so := parseOpts(opts...)
	for qo, e := range ops {
		if reflect.TypeOf(qo) == reflect.TypeOf(p.Query()) {
			return e(ctx, p)
		}
	}
	if so.defaultExecute != nil {
		return so.defaultExecute(ctx, p)
	}
	if so.panic {
		panic(fmt.Sprintf("%T not supported for model %s", p.Query(), p.Model().Type().Name()))
	}
	return nil
}

// SwitchOpt implements the options pattern for Switch.
type SwitchOpt func(s *switchOpts)

// SwitchWithoutPanic causes Switch to avoid panicking and return a nil error if an unsupported query is provided.
func SwitchWithoutPanic() SwitchOpt {
	return func(so *switchOpts) {
		so.panic = false
	}
}

// SwitchWithDefault causes Switch to use the provided Execute if no other Execute is found within the
// Ops themselves.
func SwitchWithDefault(q Execute) SwitchOpt {
	return func(so *switchOpts) {
		so.panic = false
		so.defaultExecute = q
	}
}

func parseOpts(opts ...SwitchOpt) *switchOpts {
	so := &switchOpts{panic: true}
	for _, o := range opts {
		o(so)
	}
	return so
}

type switchOpts struct {
	panic          bool
	defaultExecute Execute
}
