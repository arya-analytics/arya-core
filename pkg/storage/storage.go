package storage

// |||| ENGINE CONFIG ||||

type EngineConfig map[EngineRole]BaseEngine

func (ec EngineConfig) retrieve(r EngineRole) BaseEngine {
	return ec[r]
}

func (ec EngineConfig) mdEngine() MDEngine {
	return ec.retrieve(EngineRoleMD).(MDEngine)
}

func (ec EngineConfig) objEngine() ObjectEngine {
	return ec.retrieve(EngineRoleObject).(ObjectEngine)
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
	return s.pooler.Retrieve(s.cfg.retrieve(r))
}
