package internal

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/google/uuid"
)

type Adapter interface {
	ID() uuid.UUID
	DemandCap() int
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
	NewAdapter() (Adapter, error)
	IsAdapter(a Adapter) bool
}

// || META DATA ||

// EngineMD or the Metadata Engine is responsible for storing lightweight,
// strongly consistent data across the cluster.
type EngineMD interface {
	Engine
	query.Assemble
	Exec(ctx context.Context, p *query.Pack) error
	NewTasks(opts ...tasks.ScheduleOpt) (tasks.Schedule, error)
}

// || OBJECT ||

// EngineObject is responsible for storing bulktelem data to node localstorage data storage.
type EngineObject interface {
	Engine
	query.AssembleCreate
	query.AssembleRetrieve
	query.AssembleDelete
	query.AssembleMigrate
	Exec(ctx context.Context, p *query.Pack) error
}

// || CACHE ||

// EngineCache is responsible for storing and serving lightweight,
// ephemeral data at high speeds.
type EngineCache interface {
	Engine
	// NewTSRetrieve opens a new QueryCacheTSRetrieve.
	NewTSRetrieve() QueryCacheTSRetrieve
	// NewTSCreate opens a new QueryCacheTSCreate.
	NewTSCreate() QueryCacheTSCreate
}

// |||| QUERY ||||

type Query interface {
	Exec(ctx context.Context) error
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
