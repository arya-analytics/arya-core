// Package query holds utilities for assembling/writing queries, transporting them through arya's data layers,
// and executing them.
//
// The foundation of this package lies in separating writing queries and executing them, to allow for patterns
// like mediators and chains of responsibility to execute queries without needing to provide the facilities for
// writing them.
//
// It supplies the following query 'writers' (types that implement the Query interface):
//		Create, Update, Retrieve, and Delete.
// Each writer uses an ORM like interfaces and generates the query into a Pack.
// A Pack represents an encapsulated query that can then be transported parsed, and executed in different locations.
// See Pack for information for parsing and executing packed queries.
//
// It also supplies Assemble interfaces as well as an AssembleBase implementation for adding query assembly functionality
// your package.
//
// Finally, it provides utilities for executing queries, such as Execute and Switch. See these types for more info
// on executing a query.
//
package query

type AssembleRetrieve interface {
	NewRetrieve() *Retrieve
}

type AssembleCreate interface {
	NewCreate() *Create
}

type AssembleUpdate interface {
	NewUpdate() *Update
}

type AssembleDelete interface {
	NewDelete() *Delete
}

type Assemble interface {
	AssembleCreate
	AssembleRetrieve
	AssembleDelete
	AssembleUpdate
}

type AssembleBase struct {
	e Execute
}

func NewAssemble(e Execute) AssembleBase {
	return AssembleBase{e}
}

func (a AssembleBase) NewCreate() *Create {
	return NewCreate().BindExec(a.e)
}

func (a AssembleBase) NewRetrieve() *Retrieve {
	return NewRetrieve().BindExec(a.e)
}

func (a AssembleBase) NewUpdate() *Update {
	return NewUpdate().BindExec(a.e)
}

func (a AssembleBase) NewDelete() *Delete {
	return NewDelete().BindExec(a.e)
}
