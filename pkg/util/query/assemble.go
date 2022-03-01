package query

type Assemble interface {
	NewCreate() *Create
	NewRetrieve() *Retrieve
	NewUpdate() *Update
	NewDelete() *Delete
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
