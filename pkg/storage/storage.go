// Package storage provides an interface for interacting with a set of data stores.
//
// Engines
//
// The package relies on a set of dependency injected storage engines to read
// and write data to.
//
// Engines (Engine) can fulfill one of three roles:
//
// MDEngine - Reads and writes lightweight, strongly consistent data to storage.
// ObjectEngine - Saves bulk data to node local data storage.
// CacheEngine - High speed cache that can read and write time series data.
//
// Initialization
//
// For initializing Storage with a set of engines, see storage.New.
//
// Writing and Executing Queries
//
// For information on writing new Queries, see Storage.
package storage

// |||| STORAGE ||||

// Storage is the main interface for the storage package.
//
// Initialization
//
// Storage should not be initialized directly, and should instead be done using New.
//
// Writing and Executing Queries
//
// Storage operates as an abstract factory for creating new queries.
// See the Query specific documentation on how to write queries.
//
// Error Handling
//
// Storage returns errors of type storage.Error.
// See Error's documentation for information on different ErrorType that can be
// returned.
//
// If an unexpected error is encountered,
// will return a storage.Error with an ErrTypeUnknown. Error.Base
// can be used to access the original error.
//
// Implementing a new Engine
//
// If you're working on modifying or implementing a new Engine,
// see Engine and its sub-interfaces.
type Storage struct {
	cfg    Config
	pooler *pooler
}

// New creates a new Storage based on the provided Config.
//
// Engine Specification
//
// Storage can operate without Config.CacheEngine and/or without Config.ObjectEngine.
// However, if any queries are run that require accessing one of these data stores,
// the query will panic.
//
// Storage cannot operate without Config.MDEngine,
// as it relies on this engine to maintain consistency with other engines.
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

func (s *Storage) adapter(e Engine) (a Adapter) {
	return s.pooler.retrieve(e)
}

// |||| CONFIG ||||

// Config holds the configuration information for Storage.
// See New for information on creating Config.
type Config struct {
	MDEngine     MDEngine
	ObjectEngine ObjectEngine
	CacheEngine  CacheEngine
}
