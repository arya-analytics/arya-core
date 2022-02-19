package storage

import (
	"context"
	storage2 "github.com/arya-analytics/aryacore/pkg/models"
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
	InCatalog(model interface{}) bool
}

// || META DATA ||

// EngineMD or the Metadata Engine is responsible for storing lightweight,
// strongly consistent data across the cluster.
type EngineMD interface {
	Engine
	// NewRetrieve opens a new QueryMDRetrieve.
	NewRetrieve(a Adapter) QueryMDRetrieve
	// NewCreate opens a new QueryMDCreate.
	NewCreate(a Adapter) QueryMDCreate
	// NewDelete opens a new QueryMDDelete.
	NewDelete(a Adapter) QueryMDDelete
	// NewMigrate opens a new QueryMDMigrate.
	NewMigrate(a Adapter) QueryMDMigrate
	// NewUpdate opens a new QueryMDUpdate.
	NewUpdate(a Adapter) QueryMDUpdate
	NewTasks(a Adapter, opts ...tasks.SchedulerOpt) tasks.Scheduler
}

// || OBJECT ||

// EngineObject is responsible for storing bulk data to node localstorage data storage.
type EngineObject interface {
	Engine
	// NewRetrieve opens a new QueryObjectRetrieve.
	NewRetrieve(a Adapter) QueryObjectRetrieve
	// NewCreate opens a new QueryObjectCreate.
	NewCreate(a Adapter) QueryObjectCreate
	// NewDelete opens a new QueryObjectDelete.
	NewDelete(a Adapter) QueryObjectDelete
	// NewMigrate opens a new QueryObjectMigrate.
	NewMigrate(a Adapter) QueryObjectMigrate
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

// QueryMDRetrieve is for retrieving items from metadata storage.
type QueryMDRetrieve interface {
	QueryMDBase
	Model(model interface{}) QueryMDRetrieve
	WherePK(pk interface{}) QueryMDRetrieve
	WherePKs(pks interface{}) QueryMDRetrieve
	Relation(rel string, fields ...string) QueryMDRetrieve
	Fields(fields ...string) QueryMDRetrieve
	WhereFields(flds storage2.Fields) QueryMDRetrieve
	Count(ctx context.Context) (int, error)
}

// QueryMDCreate is for creating items in metadata storage.
type QueryMDCreate interface {
	QueryMDBase
	Model(model interface{}) QueryMDCreate
}

// QueryMDUpdate is for updating items in metadata storage.
type QueryMDUpdate interface {
	QueryMDBase
	Model(model interface{}) QueryMDUpdate
	Where(query string, args ...interface{}) QueryMDUpdate
	WherePK(pk interface{}) QueryMDUpdate
}

// QueryMDDelete is for deleting items in metadata storage.
type QueryMDDelete interface {
	QueryMDBase
	Model(model interface{}) QueryMDDelete
	WherePK(pk interface{}) QueryMDDelete
	WherePKs(pks interface{}) QueryMDDelete
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

type QueryObjectCreate interface {
	QueryObjectBase
	Model(model interface{}) QueryObjectCreate
}

type QueryObjectRetrieve interface {
	QueryObjectBase
	Model(model interface{}) QueryObjectRetrieve
	WherePK(pk interface{}) QueryObjectRetrieve
	WherePKs(pks interface{}) QueryObjectRetrieve
}

type QueryObjectDelete interface {
	QueryObjectBase
	Model(model interface{}) QueryObjectDelete
	WherePK(pk interface{}) QueryObjectDelete
	WherePKs(pks interface{}) QueryObjectDelete
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
