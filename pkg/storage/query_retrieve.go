package storage

import "context"

type retrieveQuery struct {
	baseQuery
	_mdQuery MDRetrieveQuery
}

func newRetrieve(s *Storage) *retrieveQuery {
	r := &retrieveQuery{}
	r.init(s)
	return r
}

func (r *retrieveQuery) mdQuery() MDRetrieveQuery {
	if r._mdQuery == nil {
		r._mdQuery = r.mdEngine.NewRetrieve(r.storage.adapter(EngineRoleMD))
	}
	return r._mdQuery
}

func (r *retrieveQuery) setMDQuery(q MDRetrieveQuery) {
	r._mdQuery = q
}

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
