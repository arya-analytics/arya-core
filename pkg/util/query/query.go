package query

import (
	"context"
	"fmt"
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

func (q *Pack) bindModel(m interface{}) {
	q.model = model.NewReflect(m)
}

func (q *Pack) Model() *model.Reflect {
	return q.model
}

func (q *Pack) Query() Query {
	return q.query
}

// |||| SWITCH ||||

type Ops struct {
	Create   Execute
	Retrieve Execute
	Delete   Execute
	Update   Execute
}

func Switch(ctx context.Context, p *Pack, ops Ops) error {
	switch p.Query().(type) {
	case *Create:
		return ops.Create(ctx, p)
	case *Retrieve:
		return ops.Retrieve(ctx, p)
	case *Update:
		return ops.Update(ctx, p)
	case *Delete:
		return ops.Delete(ctx, p)
	default:
		panic(fmt.Sprintf("%T not supported for model %s", p.Query(), p.Model().Type().Name()))
	}
}
