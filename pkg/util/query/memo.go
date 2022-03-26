package query

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Memo struct {
	into *model.Reflect
}

func NewMemo(into *model.Reflect) *Memo {
	return &Memo{into: into}
}

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

func (m *Memo) Add(resRfl *model.Reflect) {
	resRfl.ForEach(func(rfl *model.Reflect, i int) {
		m.into.ChainAppend(rfl)
	})
}
