package tsquery

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type AssembleRetrieve interface {
	NewRetrieve() *Retrieve
}

type AssembleCreate interface {
	NewCreate() *Create
}

type Assemble interface {
	AssembleCreate
	AssembleRetrieve
	query.AssembleExec
}

type AssembleBase struct {
	e query.Execute
}

func NewAssemble(e query.Execute) AssembleBase {
	return AssembleBase{e: e}
}

func (a AssembleBase) Exec(ctx context.Context, p *query.Pack) error {
	return a.e(ctx, p)
}

func (a AssembleBase) NewRetrieve() *Retrieve {
	return NewRetrieve().BindExec(a.e)
}

func (a AssembleBase) NewCreate() *Create {
	return NewCreate().BindExec(a.e)
}
