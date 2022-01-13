package storage

import "context"

type retrieve struct {
	base
	_mdQuery MetaDataRetrieve
}

func newRetrieve(s *Storage) *retrieve {
	r := &retrieve{}
	r.init(s)
	return r
}

func (r *retrieve) mdQuery() MetaDataRetrieve {
	if r._mdQuery == nil {
		r._mdQuery = r.mdEngine.NewRetrieve(r.storage.adapter(EngineRoleMetaData))
	}
	return r._mdQuery
}

func (r *retrieve) setMDQuery(q MetaDataRetrieve) {
	r._mdQuery = q
}

func (r *retrieve) WhereID(id interface{}) *retrieve {
	r.setMDQuery(r.mdQuery().WhereID(id))
	return r
}

func (r *retrieve) Model(model interface{}) *retrieve {
	r.setMDQuery(r.mdQuery().Model(model))
	return r
}

func (r *retrieve) Exec(ctx context.Context) error {
	return r.mdQuery().Exec(ctx)
}
