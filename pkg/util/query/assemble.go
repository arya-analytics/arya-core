package query

type Assemble struct {
	e Execute
}

func NewAssemble(e Execute) Assemble {
	return Assemble{e}
}

func (a Assemble) NewCreate() *Create {
	return NewCreate().BindExec(a.e)
}

func (a Assemble) NewRetrieve() *Retrieve {
	return NewRetrieve().BindExec(a.e)
}

func (a Assemble) NewUpdate() *Update {
	return NewUpdate().BindExec(a.e)
}

func (a Assemble) NewDelete() *Delete {
	return NewDelete().BindExec(a.e)
}
