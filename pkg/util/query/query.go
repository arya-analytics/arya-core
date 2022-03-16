// Package query holds utilities for assembling/writing queries, transporting them through arya's data layers,
// and executing them.
//
// The foundation of this package lies in separating writing queries and executing them, to allow for patterns
// like mediators and chains of responsibility to execute queries without needing to provide the facilities for
// writing them.
//
// It supplies the following query 'writers' (types that implement the Query interface):
// Create, Update, Retrieve, and Delete.
// Each writer uses an ORM like interface and 'packs' the query into a Pack.
// A Pack represents an encapsulated query that can then be transported parsed, and executed in different locations.
// See Pack for information for parsing and executing packed queries.
//
// It also supplies Assemble interfaces as well as an AssembleBase implementation for adding query assembly functionality
// to your package.
//
// Finally, it provides utilities for executing queries, such as Execute and Switch. See these types for more info
// on executing a query.
//
package query

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| QUERY ||||

// Query is a general interface for a type that can be used to write queries. (examples within this package are
// Create, Update, Retrieve, Delete).
//
// Pack packs the query into a Pack. (I know, so self-explanatory!)
//
// Exec executes the query. A query should also provide a utility for binding Execute which can be used to execute
// the query (this method should be named BindExec).For an example, see Retrieve.
// It's typical for a Query to call Execute internally when the caller calls Exec.
//
type Query interface {
	Pack() *Pack
	Exec(ctx context.Context) error
}

// |||| PACK ||||

// Pack is a representation of a query as a struct. It stores the model, variant, and options for a query.
// A Pack can be easily transported from where it's assembled to where it needs to be executed.
//
// Pack should generally not be instantiated directly, and should instead be created by using a Query such as
// Create.
//
type Pack struct {
	query Query
	model *model.Reflect
	opts  opts
}

func NewPack(q Query) *Pack {
	return &Pack{query: q, opts: map[optKey]interface{}{}}
}

func (p *Pack) bindModel(m interface{}) {
	switch m.(type) {
	case *model.Reflect:
		p.model = m.(*model.Reflect)
	default:
		p.model = model.NewReflect(m)
	}
}

// Model returns the packed query's model.
func (p *Pack) Model() *model.Reflect {
	return p.model
}

// Query returns the underlying Query the pack was built off of.
func (p *Pack) Query() Query {
	return p.query
}

// String stringifies information about the Query
func (p *Pack) String() string {
	var count = 1
	if p.Model().IsChain() {
		count = p.Model().ChainValue().Len()
	}
	return fmt.Sprintf("Variant: %T, Model: %s, Count: %v, Opts: %s", p.query, p.model.Type(), count, p.opts)
}
