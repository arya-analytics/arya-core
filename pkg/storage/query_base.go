package storage

type baseQuery struct {
	storage      *Storage
	mdEngine     MDEngine
	_mdBaseQuery MDBaseQuery
}

func (b *baseQuery) baseInit(s *Storage) {
	b.storage = s
	b.mdEngine = b.storage.cfg.mdEngine()
}

// |||| QUERY BINDING ||||

func (b *baseQuery) baseMDAdapter() Adapter {
	return b.storage.adapter(EngineRoleMD)
}

func (b *baseQuery) baseMDQuery() MDBaseQuery {
	return b._mdBaseQuery
}

func (b *baseQuery) baseSetMDQuery(q MDBaseQuery) {
	b._mdBaseQuery = q
}

func (b *baseQuery) baseObjectAdapter() Adapter {
	return b.storage.adapter(EngineRoleObject)
}
