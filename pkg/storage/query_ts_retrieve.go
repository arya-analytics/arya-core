package storage

type tsRetrieveQuery struct {
	baseQuery
}

func newTSRetrieve(s *Storage) *tsRetrieveQuery {
	tr := &tsRetrieveQuery{}
	tr.baseInit(s)
	return tr
}
