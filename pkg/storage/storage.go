// Package storage provides an interface for interacting with a set of data stores.
// The package relies on a set of dependency injected storage engines to read
// and write data to.
//
// Storage engines can fulfill one of three roles:
//		EngineRoleMD
//
//
//
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

func (s *Storage) adapter(e BaseEngine) (a Adapter) {
	return s.pooler.retrieve(e)
}

// |||| CONFIG ||||

type Config struct {
	MDEngine     MDEngine
	ObjectEngine ObjectEngine
	CacheEngine  CacheEngine
}
