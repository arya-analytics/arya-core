package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// AssembleTSRetrieve assembles a query for retrieving time series data.
type AssembleTSRetrieve interface {
	NewTSRetrieve() *TSRetrieve
}

// AssembleTSCreate assembles a query for creating time series data.
type AssembleTSCreate interface {
	NewTSCreate() *TSCreate
}

// AssembleTS composes the above interfaces into a single 'assembler.'
// Implementing this interface is ideal for types that can perform all the above operations on a data store.
type AssembleTS interface {
	AssembleTSCreate
	AssembleTSRetrieve
	query.AssembleExec
}

// AssembleTSBase is a base implementation of the AssembleTS interface.
// To create a new AssembleTSBase, call NewAssembleTS.
type AssembleTSBase struct {
	e query.Execute
}

// NewAssembleTS creates a new AssembleTSBase that will run queries against the given query.Execute implementation.
func NewAssembleTS(e query.Execute) AssembleTSBase {
	return AssembleTSBase{e: e}
}

// Exec implements the query.AssembleExec interface.
func (a AssembleTSBase) Exec(ctx context.Context, p *query.Pack) error {
	return a.e(ctx, p)
}

// NewTSRetrieve creates implements the AssembleTSRetrieve interface.
func (a AssembleTSBase) NewTSRetrieve() *TSRetrieve {
	return NewTSRetrieve().BindExec(a.e)
}

// NewTSCreate creates implements the AssembleTSCreate interface.
func (a AssembleTSBase) NewTSCreate() *TSCreate {
	return NewTSCreate().BindExec(a.e)
}
