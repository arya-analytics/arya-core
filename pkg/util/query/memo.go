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
func NewMemo(into interface{}) *Memo {
	return &Memo{into: model.NewReflect(into)}
}

// Exec implements query.Execute, and will attempt to satisfy the query using memoized results.
func (m *Memo) Exec(ctx context.Context, p *Pack) error {
	pkc, pkcOk := RetrievePKOpt(p)
	if !pkcOk {
		return NewSimpleError(ErrorTypeItemNotFound, errors.New("item not found in memo"))
	}
	for _, pk := range pkc {
		v, vOk := m.into.ValueByPK(pk)
		if !vOk {
			return NewSimpleError(ErrorTypeItemNotFound, errors.New("item not found in memo"))
		}
		if p.Model().IsStruct() {
			p.Model().Set(v)
		} else {
			p.Model().ChainAppend(v)
		}
	}
	return nil
}

// Add adds a query result to the memo.
func (m *Memo) Add(resRfl *model.Reflect) {
	resRfl.ForEach(func(rfl *model.Reflect, i int) {
		if m.into.IsStruct() {
			m.into.Set(rfl)
		} else {
			m.into.ChainAppend(rfl)
		}
	})
}
