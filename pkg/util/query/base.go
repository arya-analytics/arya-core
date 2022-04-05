package query

import (
	"context"
)

// Base is the base type for all query types.
type Base struct {
	e  Execute
	_q *Pack
}

// || CONSTRUCTOR ||

// Init initializes the query.
func (b *Base) Init(q Query) {
	b._q = NewPack(q)
}

// || PACK ||

// Pack packs the query.
func (b *Base) Pack() *Pack {
	return b._q
}

// || MODEL ||

// Model binds the model for the query.
func (b *Base) Model(m interface{}) {
	b.Pack().bindModel(m)
}

// || EXECUTION ||

// BindExec binds Execute implementation for the query.
func (b *Base) BindExec(e Execute) {
	b.e = e
}

// Exec executes the query. It returns any errors encountered during execution.
// Results from the query will be bound to the parameters passed in during assembly.
func (b *Base) Exec(ctx context.Context) error {
	p := b.Pack()
	memo, ok := MemoOpt(p)
	if ok {
		if err := memo.Exec(ctx, p); err == nil {
			return nil
		}
	}
	if b.e == nil {
		panic("query execute not bound")
	}
	err := b.e(ctx, p)
	if ok && err == nil {
		memo.Add(p.Model())
	}
	return err

}
