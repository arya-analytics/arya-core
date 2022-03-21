package query

type StreamRetrieve struct {
	where
}

func (sr *StreamRetrieve) Model(m interface{}) *StreamRetrieve {
	sr.baseModel(m)
	return sr
}

func (sr *StreamRetrieve) WherePK(pk interface{}) *StreamRetrieve {
	sr.wherePK(pk)
	return sr
}

func (sr *StreamRetrieve) WherePKs(pks interface{}) *StreamRetrieve {
	sr.wherePKs(pks)
	return sr
}

func (sr *StreamRetrieve) BindExec(e Execute) *StreamRetrieve {
	sr.baseBindExec(e)
	return sr
}
