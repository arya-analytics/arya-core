package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/google/uuid"
)

type Adapter interface {
	ID() uuid.UUID
}

// |||| ENGINE ||||

// Engine is a set of general interfaces that each engine variant must meet.
//
// Assigning Telem Responsibility
//
// Each engine variant is responsible for storing specific data types.
// These responsibilities are assigned in the model struct using the storage.re key.
// If no responsibility is assigned, EngineMD is assumed responsible.
type Engine interface {
	NewAdapter() Adapter
	IsAdapter(a Adapter) bool
	ShouldHandle(model interface{}, flds ...string) bool
}

// || META DATA ||

// EngineMD or the Metadata Engine is responsible for storing lightweight,
// strongly consistent data across the cluster.
type EngineMD interface {
	Engine
	query.Assemble
	NewMigrate() QueryMDMigrate
	NewTasks(opts ...tasks.ScheduleOpt) tasks.Schedule
}

// || OBJECT ||

// EngineObject is responsible for storing chanchunk data to node localstorage data storage.
type EngineObject interface {
	Engine
	query.AssembleCreate
	query.AssembleRetrieve
	query.AssembleDelete
	NewMigrate() QueryObjectMigrate
}

// || CACHE ||

// EngineCache is responsible for storing and serving lightweight,
// ephemeral data at high speeds.
type EngineCache interface {
	Engine
	// NewTSRetrieve opens a new QueryCacheTSRetrieve.
	NewTSRetrieve(a Adapter) QueryCacheTSRetrieve
	// NewTSCreate opens a new QueryCacheTSCreate.
	NewTSCreate(a Adapter) QueryCacheTSCreate
}

// |||| QUERY ||||

type Query interface {
	Exec(ctx context.Context) error
}

// || META DATA ||

type QueryMDBase interface {
	Query
}

// QueryMDMigrate applies migration changes to metadata storage.
type QueryMDMigrate interface {
	QueryMDBase
	Verify(ctx context.Context) error
}

// || OBJECT ||

type QueryObjectBase interface {
	Query
}

type QueryObjectMigrate interface {
	QueryObjectBase
	Verify(ctx context.Context) error
}

// || TS CACHE ||

type QueryCacheBase interface {
	Query
}

type QueryCacheCreate interface {
	QueryCacheBase
	Model(model interface{}) QueryCacheCreate
}

type QueryCacheTSRetrieve interface {
	QueryCacheBase
	SeriesExists(ctx context.Context, pk interface{}) (bool, error)
	Model(model interface{}) QueryCacheTSRetrieve
	WherePK(pk interface{}) QueryCacheTSRetrieve
	WherePKs(pks interface{}) QueryCacheTSRetrieve
	AllTimeRange() QueryCacheTSRetrieve
	WhereTimeRange(fromTS int64, toTS int64) QueryCacheTSRetrieve
}

type QueryCacheTSCreate interface {
	QueryCacheBase
	Model(model interface{}) QueryCacheTSCreate
	Series() QueryCacheTSCreate
	Sample() QueryCacheTSCreate
}
