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
// EngineObject - Saves bulktelem data to node localstorage data storage.
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
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
	query.Assemble
	Exec(ctx context.Context, p *query.Pack) error
	NewTSRetrieve() *QueryTSRetrieve
	NewTSCreate() *QueryTSCreate
	NewMigrate() *QueryMigrate
	AddQueryHook(hook QueryHook)
	Start(ctx context.Context, opts ...tasks.ScheduleOpt)
	Stop()
	Errors() chan error
}

type storage struct {
	query.AssembleBase
	cfg        Config
	tasks      tasks.Schedule
	queryHooks []QueryHook
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
	s := &storage{cfg: cfg}
	s.AssembleBase = query.NewAssemble(s.Exec)
	return s
}

// Exec implements query.Execute
func (s *storage) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		Create:   newDef(s).exec,
		Retrieve: newDef(s).exec,
		Delete:   newDef(s).exec,
		Update:   newUpdate(s).exec,
	})
}

// NewMigrate opens a new QueryMigrate.
func (s *storage) NewMigrate() *QueryMigrate {
	return newMigrate(s)
}

// NewTSRetrieve opens a new QueryTSRetrieve.
func (s *storage) NewTSRetrieve() *QueryTSRetrieve {
	return newTSRetrieve(s)
}

// NewTSCreate opens a new QueryTSCreate.
func (s *storage) NewTSCreate() *QueryTSCreate {
	return newTSCreate(s)
}

// Start starts
func (s *storage) Start(ctx context.Context, opts ...tasks.ScheduleOpt) {
	s.tasks = tasks.NewScheduleBatch(s.cfg.EngineMD.NewTasks(opts...))
	s.tasks.Start(ctx)
}

func (s *storage) Stop() {
	s.tasks.Stop()
	s.tasks = nil
}

func (s *storage) Errors() chan error {
	return s.tasks.Errors()
}

func (s *storage) AddQueryHook(hook QueryHook) {
	s.queryHooks = append(s.queryHooks, hook)
}

// |||| CONFIG ||||

// Config holds the configuration information for Storage.
// See New for information on creating Config.
type Config struct {
	EngineMD     EngineMD
	EngineObject EngineObject
	EngineCache  EngineCache
}
