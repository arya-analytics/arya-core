package query

import "context"

// AssembleRetrieve is used to create queries for retrieving items from a data store.
type AssembleRetrieve interface {
	NewRetrieve() *Retrieve
}

// AssembleCreate is used to create queries for creating items in a data store.
type AssembleCreate interface {
	NewCreate() *Create
}

// AssembleUpdate is used to create queries for updating items in a data store.
type AssembleUpdate interface {
	NewUpdate() *Update
}

// AssembleDelete is used to create queries for deleting items from a data store.
type AssembleDelete interface {
	NewDelete() *Delete
}

// AssembleMigrate is used to create queries for migrating items in a data store.
type AssembleMigrate interface {
	NewMigrate() *Migrate
}

// AssembleExec is an interface that implements query.Execute, which can be used to execute generic packed queries.
type AssembleExec interface {
	Exec(ctx context.Context, p *Pack) error
}

// Assemble composes the above interfaces into a single query 'assembler'.
// Implementing this interface is ideal for types that can perform all of the above operations on a data store.
type Assemble interface {
	AssembleCreate
	AssembleRetrieve
	AssembleDelete
	AssembleUpdate
	AssembleMigrate
	AssembleExec
}

// AssembleBase is a Base implementation of the Assemble interface.
// To create a new AssembleBase, call NewAssembleBase.
type AssembleBase struct {
	e Execute
}

// NewAssemble initializes a new AssembleBase that will run queries against the given Execute implementation.
func NewAssemble(e Execute) Assemble {
	return AssembleBase{e}
}

// Exec implements the AssembleExec interface.
func (a AssembleBase) Exec(ctx context.Context, p *Pack) error {
	return a.e(ctx, p)
}

// NewCreate implements the AssembleCreate interface.
func (a AssembleBase) NewCreate() *Create {
	return NewCreate().BindExec(a.e)
}

// NewRetrieve implements the AssembleRetrieve interface.
func (a AssembleBase) NewRetrieve() *Retrieve {
	return NewRetrieve().BindExec(a.e)
}

// NewUpdate implements the AssembleUpdate interface.
func (a AssembleBase) NewUpdate() *Update {
	return NewUpdate().BindExec(a.e)
}

// NewDelete implements the AssembleDelete interface.
func (a AssembleBase) NewDelete() *Delete {
	return NewDelete().BindExec(a.e)
}

// NewMigrate implements the AssembleMigrate interface.
func (a AssembleBase) NewMigrate() *Migrate {
	return NewMigrate().BindExec(a.e)
}
