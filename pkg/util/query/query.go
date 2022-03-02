package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| QUERY ||||

// Query is a general interface for a type that can be used to write queries. (examples within this package are
// Create, Update, Retrieve, Delete).
//
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

func (q *Pack) bindModel(m interface{}) {
	switch m.(type) {
	case *model.Reflect:
		q.model = m.(*model.Reflect)
	default:
		q.model = model.NewReflect(m)
	}
}

// Model returns the packed query's model.
func (q *Pack) Model() *model.Reflect {
	return q.model
}

// Query returns the underlying Query the pack was built off of.
func (q *Pack) Query() Query {
	return q.query
}
