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
	if ok {
		memo.Add(p.Model())
	}
	return err

}
