package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| QUERY ||||

type Query interface {
	baseBindExec(e Execute)
	Pack() *Pack
	Exec(ctx context.Context) error
}

type Pack struct {
	query Query
	model *model.Reflect
	opts  opts
}

func NewPack(q Query) *Pack {
	return &Pack{query: q, opts: map[optKey]interface{}{}}
}

func (q *Pack) BindModel(m interface{}) {
	q.model = model.NewReflect(m)
}

func (q *Pack) Model() *model.Reflect {
	return q.model
}

func (q *Pack) Query() Query {
	return q.query
}
