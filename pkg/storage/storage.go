package storage

// |||| CONFIG ||||

type Config map[EngineRole]BaseEngine

func (ec Config) retrieve(r EngineRole) BaseEngine {
	return ec[r]
}

func (ec Config) mdEngine() MDEngine {
	return ec.retrieve(EngineRoleMD).(MDEngine)
}

func (ec Config) objEngine() ObjectEngine {
	return ec.retrieve(EngineRoleObject).(ObjectEngine)
}

func (ec Config) cacheEngine() CacheEngine {
	return ec.retrieve(EngineRoleCache).(CacheEngine)
}

// |||| STORAGE ||||

type Storage struct {
	cfg    Config
	pooler *pooler
}

func New(cfg Config) *Storage {
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

func (s *Storage) NewTSCreate() *tsCreateQuery {
	return newTSCreate(s)
}

func (s *Storage) adapter(r EngineRole) (a Adapter) {
	e := s.cfg.retrieve(r)
	return s.pooler.Retrieve(e)
}
