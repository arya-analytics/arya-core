package storage

type baseQuery struct {
	storage       *Storage
	mdEngine      MDEngine
	_baseMDQuery  MDBaseQuery
	objEngine     ObjectEngine
	_baseObjQuery ObjectBaseQuery
}

func (b *baseQuery) baseInit(s *Storage) {
	b.storage = s
	b.mdEngine = b.storage.cfg.mdEngine()
	b.objEngine = b.storage.cfg.objEngine()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (b *baseQuery) baseMDAdapter() Adapter {
	return b.storage.adapter(EngineRoleMD)
}

func (b *baseQuery) baseMDQuery() MDBaseQuery {
	return b._baseMDQuery
}

func (b *baseQuery) baseSetMDQuery(q MDBaseQuery) {
	b._baseMDQuery = q
}

// || OBJECT ||

func (b *baseQuery) baseObjAdapter() Adapter {
	return b.storage.adapter(EngineRoleObject)
}

func (b *baseQuery) baseObjQuery() ObjectBaseQuery {
	return b._baseObjQuery
}

func (b *baseQuery) baseSetObjQuery(q ObjectBaseQuery) {
	b._baseObjQuery = q
}
