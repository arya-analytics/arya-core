package storage

// |||| ENGINE CONFIG ||||

type EngineConfig map[EngineRole]BaseEngine

func (ec EngineConfig) retrieve(r EngineRole) (BaseEngine, bool) {
	e, ok := ec[r]
	return e, ok
}

func (ec EngineConfig) mdEngine() MDEngine {
	e, ok := ec.retrieve(EngineRoleMD)
	if !ok {
		return nil
	}
	return e.(MDEngine)
}

func (ec EngineConfig) objEngine() ObjectEngine {
	e, ok := ec.retrieve(EngineRoleObject)
	if !ok {
		return nil
	}
	return e.(ObjectEngine)
}

func (ec EngineConfig) cacheEngine() CacheEngine {
	e, ok := ec.retrieve(EngineRoleCache)
	if !ok {
		return nil
	}
	return e.(CacheEngine)
}

// |||| STORAGE ||||

type Storage struct {
	cfg    EngineConfig
	pooler *pooler
}

func New(cfg EngineConfig) *Storage {
	return &Storage{
		cfg:    cfg,
		pooler: newPooler(),
	}
}

func (s *Storage) NewMigrate() *migrateQuery {
	return newMigrate(s)
}

func (s *Storage) NewRetrieve() *retrieveQuery {
	return newRetrieve(s)
}

func (s *Storage) NewCreate() *createQuery {
	return newCreate(s)
}

func (s *Storage) NewDelete() *deleteQuery {
	return newDelete(s)
}

func (s *Storage) NewUpdate() *updateQuery {
	return newUpdate(s)
}

func (s *Storage) NewTSRetrieve() *tsRetrieveQuery {
	return newTSRetrieve(s)
}

func (s *Storage) adapter(r EngineRole) (a Adapter) {
	e, ok := s.cfg.retrieve(r)
	if !ok {
		panic("tried to retrieve a non-existent engine")
	}
	return s.pooler.Retrieve(e)
}
