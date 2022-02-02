package storage

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

// NewMigrate opens a new MigrateQuery.
func (s *Storage) NewMigrate() *MigrateQuery {
	return newMigrate(s)
}

// NewRetrieve opens a new RetrieveQuery.
func (s *Storage) NewRetrieve() *RetrieveQuery {
	return newRetrieve(s)
}

// NewCreate opens a new CreateQuery.
func (s *Storage) NewCreate() *CreateQuery {
	return newCreate(s)
}

// NewDelete opens a new DeleteQuery.
func (s *Storage) NewDelete() *DeleteQuery {
	return newDelete(s)
}

// NewUpdate opens a new UpdateQuery.
func (s *Storage) NewUpdate() *UpdateQuery {
	return newUpdate(s)
}

// NewTSRetrieve opens a new TSRetrieveQuery.
func (s *Storage) NewTSRetrieve() *TSRetrieveQuery {
	return newTSRetrieve(s)
}

// NewTSCreate opens a new TSCreateQuery.
func (s *Storage) NewTSCreate() *TSCreateQuery {
	return newTSCreate(s)
}

func (s *Storage) adapter(r EngineRole) (a Adapter) {
	e := s.cfg.retrieve(r)
	return s.pooler.retrieve(e)
}

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
