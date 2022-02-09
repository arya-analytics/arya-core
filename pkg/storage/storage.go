// Package storage provides an interface for interacting with a set of data stores.
//
// Storage architecture documentation can be found at:
// https://arya-analytics.atlassian.net/wiki/spaces/AA/pages/3212201/00+-+Storage+Layer
// High level data architecture documentation can be found at:
// https://arya-analytics.atlassian.net/wiki/spaces/AA/pages/819257/00+-+Arya+Core#5.2---Data-Architecture
//
// Engines
//
// The package relies on a set of dependency injected storage engines to read
// and write data to.
//
// Engines (Engine) can fulfill one of three roles:
//
// EngineMD - Reads and writes lightweight, strongly consistent data to storage.
// EngineObject - Saves bulk data to node local data storage.
// EngineCache - High speed cache that can read and write time series data.
//
// Initialization
//
// For initializing Storage with a set of engines, see storage.New.
//
// Writing and Executing Queries
//
// For information on writing new Queries, see Storage.
package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

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
// will return a storage.Error with an ErrorTypeUnknown. Error.Base
// can be used to access the original error.
//
// Implementing a new Engine
//
// If you're working on modifying or implementing a new Engine,
// see Engine and its sub-interfaces.
type Storage interface {
	NewCreate() *QueryCreate
	NewRetrieve() *QueryRetrieve
	NewUpdate() *QueryUpdate
	NewDelete() *QueryDelete
	NewTSRetrieve() *QueryTSRetrieve
	NewTSCreate() *QueryTSCreate
	NewMigrate() *QueryMigrate
	NewTasks(opts ...tasks.SchedulerOpt) tasks.Scheduler
	config() Config
	adapter(e Engine) Adapter
}

type storage struct {
	cfg           Config
	pooler        *pooler
	taskScheduler tasks.Scheduler
}

// New creates a new Storage based on the provided Config.
//
// Engine Specification
//
// Storage can operate without Config.EngineCache and/or without Config.EngineObject.
// However, if any queries are run that require accessing one of these data stores,
// the query will panic.
//
// Storage cannot operate without Config.EngineMD,
// as it relies on this engine to maintain consistency with other engines.
func New(cfg Config) Storage {
	return &storage{cfg: cfg, pooler: newPooler()}
}

// NewMigrate opens a new QueryMigrate.
func (s *storage) NewMigrate() *QueryMigrate {
	return newMigrate(s)
}

// NewRetrieve opens a new QueryRetrieve.
func (s *storage) NewRetrieve() *QueryRetrieve {
	return newRetrieve(s)
}

// NewCreate opens a new QueryCreate.
func (s *storage) NewCreate() *QueryCreate {
	return newCreate(s)
}

// NewDelete opens a new QueryDelete.
func (s *storage) NewDelete() *QueryDelete {
	return newDelete(s)
}

// NewUpdate opens a new QueryUpdate.
func (s *storage) NewUpdate() *QueryUpdate {
	return newUpdate(s)
}

// NewTSRetrieve opens a new QueryTSRetrieve.
func (s *storage) NewTSRetrieve() *QueryTSRetrieve {
	return newTSRetrieve(s)
}

// NewTSCreate opens a new QueryTSCreate.
func (s *storage) NewTSCreate() *QueryTSCreate {
	return newTSCreate(s)
}

func (s *storage) NewTasks(opts ...tasks.SchedulerOpt) tasks.Scheduler {
	return tasks.NewBatchScheduler(
		s.cfg.EngineMD.NewTasks(s.adapter(s.cfg.EngineMD), opts...),
	)
}

func (s *storage) adapter(e Engine) (a Adapter) {
	return s.pooler.retrieve(e)
}

func (s *storage) config() Config {
	return s.cfg
}

// |||| CONFIG ||||

// Config holds the configuration information for Storage.
// See New for information on creating Config.
type Config struct {
	EngineMD     EngineMD
	EngineObject EngineObject
	EngineCache  EngineCache
}
