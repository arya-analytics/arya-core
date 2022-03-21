package query

import (
	"context"
)

type Base struct {
	e  Execute
	_q *Pack
}

// || CONSTRUCTOR ||

func (b *Base) Init(q Query) {
	b._q = NewPack(q)
}

// || PACK ||

func (b *Base) Pack() *Pack {
	return b._q
}

// || MODEL ||

func (b *Base) Model(m interface{}) {
	b.Pack().bindModel(m)
}

// || EXECUTION ||

func (b *Base) BindExec(e Execute) {
	b.e = e
}

func (b *Base) Exec(ctx context.Context) error {
	p := b.Pack()
	err := b.e(ctx, p)
	return err

}
