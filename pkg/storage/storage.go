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
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
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
// will return a storage.Error with an ErrorTypeUnknown. Error.base
// can be used to access the original error.
//
// Implementing a new Engine
//
// If you're working on modifying or implementing a new Engine,
// see Engine and its sub-interfaces.
type Storage interface {
	query.Assemble
	streamq.AssembleTS
	AddQueryHook(hook query.Hook)
	ClearQueryHooks()
	Start(ctx context.Context, opts ...tasks.ScheduleOpt) error
	Stop()
	Errors() chan error
}

type storage struct {
	query.Assemble
	streamq.AssembleTS
	cfg Config
	ts  tasks.Schedule
	*query.HookRunner
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
	s := &storage{cfg: cfg, HookRunner: query.NewHookRunner()}
	s.Assemble = query.NewAssemble(s.Exec)
	s.AssembleTS = streamq.NewAssembleTS(s.Exec)
	return s
}

// Exec implements query.Execute
func (s *storage) Exec(ctx context.Context, p *query.Pack) error {
	qc := query.NewCatch(ctx, p)
	qc.Exec(s.Before)
	qc.Exec(s.cfg.EngineMD.Exec)
	qc.Exec(s.cfg.EngineObject.Exec)
	qc.Exec(s.cfg.EngineCache.Exec)
	qc.Exec(s.After)
	return qc.Error()
}

// Start starts storage internal ts.
func (s *storage) Start(ctx context.Context, opts ...tasks.ScheduleOpt) error {
	mdT, err := s.cfg.EngineMD.NewTasks(opts...)
	if err != nil {
		return err
	}
	s.ts = tasks.NewScheduleBatch(mdT)
	go s.ts.Start(ctx)
	return err
}

// Stop stops storage internal ts.
func (s *storage) Stop() {
	if s.ts != nil {
		s.ts.Stop()
		s.ts = nil
	}
}

func (s *storage) Errors() chan error {
	return s.ts.Errors()
}

// |||| CONFIG ||||

// Config holds the configuration information for Storage.
// See New for information on creating Config.
type Config struct {
	EngineMD     internal.EngineMD
	EngineObject internal.EngineObject
	EngineCache  internal.EngineCache
}
