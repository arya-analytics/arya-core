package tsquery

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type AssembleTSRetrieve interface {
	NewTSRetrieve() *Retrieve
}

type AssembleTSCreate interface {
	NewTSCreate() *Create
}

type AssembleTS interface {
	AssembleTSCreate
	AssembleTSRetrieve
	query.AssembleExec
}

type AssembleTSBase struct {
	e query.Execute
}

func NewAssemble(e query.Execute) AssembleTSBase {
	return AssembleTSBase{e: e}
}

func (a AssembleTSBase) Exec(ctx context.Context, p *query.Pack) error {
	return a.e(ctx, p)
}

func (a AssembleTSBase) NewTSRetrieve() *Retrieve {
	return NewRetrieve().BindExec(a.e)
}

func (a AssembleTSBase) NewTSCreate() *Create {
	return NewCreate().BindExec(a.e)
}
