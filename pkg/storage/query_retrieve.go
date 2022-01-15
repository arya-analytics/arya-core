package storage

import "context"

type retrieveQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s *Storage) *retrieveQuery {
	r := &retrieveQuery{}
	r.baseInit(s)
	return r
}

// |||| INTERFACE ||||

func (r *retrieveQuery) WhereID(id interface{}) *retrieveQuery {
	r.setMDQuery(r.mdQuery().WhereID(id))
	return r
}

func (r *retrieveQuery) Model(model interface{}) *retrieveQuery {
	r.setMDQuery(r.mdQuery().Model(model))
	return r
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	return r.mdQuery().Exec(ctx)
}

// |||| QUERY BINDING ||||

func (r *retrieveQuery) mdQuery() MDRetrieveQuery {
	if r.baseMDQuery() == nil {
		r.setMDQuery(r.mdEngine.NewRetrieve(r.storage.adapter(EngineRoleMD)))
	}
	return r.baseMDQuery().(MDRetrieveQuery)
}

func (r *retrieveQuery) setMDQuery(q MDRetrieveQuery) {
	r.baseSetMDQuery(q)
}