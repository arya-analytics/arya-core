package query

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// Memo memoizes the result of a query by its primary key.
// Avoid instantiating this directly, instead use NewMemo.
// The results of memoized queries are never invalidated, so use with care.
type Memo struct {
	into *model.Reflect
}

// NewMemo instantiates a new Memo that saves and looks up query results from into.
// NewMemo panics if into is not a chain.
func NewMemo(into *model.Reflect) *Memo {
	if !into.IsChain() {
		panic("memo requires a chain")
	}
	return &Memo{into: into}
}

// Exec implements query.Execute, and will attempt to satisfy the query using memoized results.
func (m *Memo) Exec(ctx context.Context, p *Pack) error {
	pkc, pkcOk := PKOpt(p)
	if !pkcOk {
		return NewSimpleError(ErrorTypeItemNotFound, errors.New("item not found in memo"))
	}
	for _, pk := range pkc {
		v, vOk := m.into.ValueByPK(pk)
		if vOk {
			if p.Model().IsStruct() {
				p.Model().PointerValue().Set(v.PointerValue())
			}
			if p.Model().IsChain() {
				p.Model().ChainAppend(v)
			}
		} else {
			return NewSimpleError(ErrorTypeItemNotFound, errors.New("item not found in memo"))
		}
	}
	return nil
}

// Add adds a query result to the memo.
func (m *Memo) Add(resRfl *model.Reflect) {
	resRfl.ForEach(func(rfl *model.Reflect, i int) {
		m.into.ChainAppend(rfl)
	})
}
