package query

import (
	"context"
)

type base struct {
	e  Execute
	_q *Pack
}

// || CONSTRUCTOR ||

func (b *base) baseInit(q Query) {
	b._q = NewPack(q)
}

// || PACK ||

func (b *base) Pack() *Pack {
	return b._q
}

// || MODEL ||

func (b *base) baseModel(m interface{}) {
	b.Pack().bindModel(m)
}

// || EXECUTION ||

func (b *base) baseBindExec(e Execute) {
	b.e = e
}

func (b *base) Exec(ctx context.Context) error {
	return b.e(ctx, b.Pack())
}
